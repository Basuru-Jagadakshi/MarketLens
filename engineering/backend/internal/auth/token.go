package auth

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ThunderID configuration, read once at startup. Call ValidateConfig
// before serving any requests — a missing THUNDER_ISSUER in particular
// must fail fast here, since jwt/v5 silently accepts a token with ANY
// issuer when WithIssuer is given an empty string.
var (
	thunderBaseURL     = os.Getenv("THUNDER_BASE_URL")
	thunderIssuer      = os.Getenv("THUNDER_ISSUER")
	thunderAudience    = os.Getenv("THUNDER_AUDIENCE")
	thunderInsecureTLS = os.Getenv("THUNDER_INSECURE_TLS") == "true"
)

const (
	// jwksRefreshCooldown limits how often a cache-miss (an unrecognized
	// kid) can trigger a real HTTP fetch, so a token with a bogus kid
	// can't be used to hammer ThunderID's JWKS endpoint on every request.
	jwksRefreshCooldown = 10 * time.Second
	// httpClientTimeout bounds outbound HTTP calls to ThunderID.
	httpClientTimeout = 10 * time.Second
	// jwtClockLeeway tolerates clock skew when validating exp/nbf.
	jwtClockLeeway = 30 * time.Second
)

// ValidateConfig reports which required environment variables are
// missing. THUNDER_INSECURE_TLS is intentionally excluded — unset
// means "use real TLS verification", which is a valid, safe default.
func ValidateConfig() error {
	var missing []string
	if thunderBaseURL == "" {
		missing = append(missing, "THUNDER_BASE_URL")
	}
	if thunderIssuer == "" {
		missing = append(missing, "THUNDER_ISSUER")
	}
	if thunderAudience == "" {
		missing = append(missing, "THUNDER_AUDIENCE")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variable(s): %s", strings.Join(missing, ", "))
	}
	return nil
}

// customClaims adds the space-delimited "scope" claim ThunderID issues
// (e.g. "openid book:create book:update") on top of the standard
// iss/sub/exp/aud claims, which jwt.RegisteredClaims already parses
// (including the string-or-array "aud" shape the JWT spec allows).
type customClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	jwksMu          sync.Mutex
	jwksCache       []jwk
	jwksLastRefresh time.Time
)

func newHTTPClient() *http.Client {
	client := &http.Client{Timeout: httpClientTimeout}

	if thunderInsecureTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return client
}

func getJWKS() ([]jwk, error) {
	jwksMu.Lock()
	defer jwksMu.Unlock()
	if jwksCache != nil {
		return jwksCache, nil
	}
	return fetchJWKSLocked()
}

// refreshJWKS re-fetches the JWKS unconditionally, replacing the cache —
// called when a token's kid isn't found in the cached set, which usually
// means ThunderID rotated its signing keys since the last fetch.
func refreshJWKS() ([]jwk, error) {
	jwksMu.Lock()
	defer jwksMu.Unlock()
	if time.Since(jwksLastRefresh) < jwksRefreshCooldown {
		return jwksCache, nil
	}
	return fetchJWKSLocked()
}

// fetchJWKSLocked does the actual HTTP fetch and cache update. Callers must
// hold jwksMu.
func fetchJWKSLocked() ([]jwk, error) {
	client := newHTTPClient()

	resp, err := client.Get(thunderBaseURL + "/oauth2/jwks")
	if err != nil {
		return nil, errors.New("cannot reach ThunderID JWKS")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ThunderID JWKS returned status %d", resp.StatusCode)
	}

	var doc struct {
		Keys []jwk `json:"keys"`
	}
	if json.NewDecoder(resp.Body).Decode(&doc) != nil || len(doc.Keys) == 0 {
		return nil, errors.New("bad JWKS document")
	}
	jwksCache = doc.Keys
	jwksLastRefresh = time.Now()
	return jwksCache, nil
}

// Build an *rsa.PublicKey from a JWK's modulus (n) and exponent (e).
func translatePublicKeyFromJWK(k jwk) (*rsa.PublicKey, error) {
	nBytes, err1 := base64.RawURLEncoding.DecodeString(k.N)
	eBytes, err2 := base64.RawURLEncoding.DecodeString(k.E)
	if err1 != nil || err2 != nil || len(eBytes) > 8 {
		return nil, errors.New("bad JWK")
	}
	padded := make([]byte, 8)
	copy(padded[8-len(eBytes):], eBytes)
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(binary.BigEndian.Uint64(padded)),
	}, nil
}

// Full verification: signature + standard claims

// jwksKeyfunc resolves the RSA public key to verify a token's signature
// with, by looking up its "kid" header in ThunderID's JWKS. An unrecognized
// kid triggers one refresh before giving up, in case ThunderID rotated its
// signing keys since the cache was last populated.
func jwksKeyfunc(t *jwt.Token) (any, error) {
	kid, _ := t.Header["kid"].(string)

	keys, err := getJWKS()
	if err != nil {
		return nil, err
	}
	if key, ok := findJWK(keys, kid); ok {
		return translatePublicKeyFromJWK(key)
	}

	keys, err = refreshJWKS()
	if err != nil {
		return nil, err
	}
	if key, ok := findJWK(keys, kid); ok {
		return translatePublicKeyFromJWK(key)
	}
	return nil, errors.New("signing key not found")
}

func findJWK(keys []jwk, kid string) (jwk, bool) {
	for _, k := range keys {
		if k.Kid == kid {
			return k, true
		}
	}
	return jwk{}, false
}

func verifyToken(raw string) (*customClaims, error) {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"RS256"}), // algorithm pinning: the server decides, never the token.
		jwt.WithIssuer(thunderIssuer),
		jwt.WithExpirationRequired(), // a token with no exp claim is rejected, not treated as never-expiring.
		jwt.WithLeeway(jwtClockLeeway),
		jwt.WithAudience(thunderAudience),
	}

	var claims customClaims
	if _, err := jwt.ParseWithClaims(raw, &claims, jwksKeyfunc, opts...); err != nil {
		return nil, err
	}
	return &claims, nil
}
