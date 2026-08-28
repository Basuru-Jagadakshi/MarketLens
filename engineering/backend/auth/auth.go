package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type jwtToken struct {
	Header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	Claims struct {
		Iss   string          `json:"iss"`
		Sub   string          `json:"sub"`
		Exp   int64           `json:"exp"`
		Scope string          `json:"scope"` // space-delimited: "openid book:create book:update"
		Aud   json.RawMessage `json:"aud"`   // spec allows a string OR an array — parsed below
	}
	SigningInput string // raw "header.payload" — the bytes the signature covers
	Signature    []byte
}

// audience returns aud as a slice whether the token used a string or an array.
func (t *jwtToken) audience() []string {
	if len(t.Claims.Aud) == 0 {
		return nil
	}
	var single string
	if json.Unmarshal(t.Claims.Aud, &single) == nil {
		return []string{single}
	}
	var many []string
	if json.Unmarshal(t.Claims.Aud, &many) == nil {
		return many
	}
	return nil
}

func parseJWT(raw string) (*jwtToken, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	var t jwtToken
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("bad token encoding")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("bad token encoding")
	}
	t.Signature, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("bad token encoding")
	}
	if json.Unmarshal(headerBytes, &t.Header) != nil || json.Unmarshal(payloadBytes, &t.Claims) != nil {
		return nil, errors.New("bad token JSON")
	}
	t.SigningInput = parts[0] + "." + parts[1]
	return &t, nil
}

// JWKS: ThunderID's public keys, fetched once and cached 

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	jwksMu    sync.Mutex
	jwksCache []jwk
)

func newHTTPClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}

	if os.Getenv("THUNDER_INSECURE_TLS") == "true" {
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

	client := newHTTPClient() 

	resp, err := client.Get(os.Getenv("THUNDER_BASE_URL") + "/oauth2/jwks")
	if err != nil {
		return nil, errors.New("cannot reach ThunderID JWKS")
	}
	defer resp.Body.Close()

	var doc struct {
		Keys []jwk `json:"keys"`
	}
	if json.NewDecoder(resp.Body).Decode(&doc) != nil {
		return nil, errors.New("bad JWKS document")
	}
	jwksCache = doc.Keys
	return jwksCache, nil
}

// Build an *rsa.PublicKey from a JWK's modulus (n) and exponent (e).
func publicKeyFromJWK(k jwk) (*rsa.PublicKey, error) {
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

func verifyToken(raw string) (*jwtToken, error) {
	t, err := parseJWT(raw)
	if err != nil {
		return nil, err
	}

	// Algorithm pinning: the server decides, never the token.
	if t.Header.Alg != "RS256" {
		return nil, errors.New("unsupported token algorithm")
	}

	keys, err := getJWKS()
	if err != nil {
		return nil, err
	}
	var key *rsa.PublicKey
	for _, k := range keys {
		if k.Kid == t.Header.Kid {
			if key, err = publicKeyFromJWK(k); err != nil {
				return nil, err
			}
			break
		}
	}
	if key == nil {
		return nil, errors.New("signing key not found")
	}

	// RS256 = RSA PKCS1v15 signature over a SHA-256 digest.
	digest := sha256.Sum256([]byte(t.SigningInput))
	if rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], t.Signature) != nil {
		return nil, errors.New("invalid token signature")
	}

	if t.Claims.Iss != os.Getenv("THUNDER_ISSUER") {
		return nil, errors.New("invalid token issuer")
	}
	if t.Claims.Exp != 0 && t.Claims.Exp < time.Now().Unix()-30 { // 30s clock skew
		return nil, errors.New("token has expired")
	}

	if audience := os.Getenv("THUNDER_AUDIENCE"); audience != "" {
		found := false
		for _, a := range t.audience() {
			if a == audience {
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("invalid token audience")
		}
	}

	return t, nil
}

// Gin middleware 

// AuthRequired: AUTHENTICATION. Verifies the JWT and stores the verified
// identity (sub) and permissions (scopes) in the request context.
// Failure → 401: we don't know who you are.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		t, err := verifyToken(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("sub", t.Claims.Sub)
		// strings.Fields splits on any whitespace and drops empties, so a
		// missing scope claim safely becomes an empty list (deny everything).
		c.Set("scopes", strings.Fields(t.Claims.Scope))
		c.Next()
	}
}

// RequireScope: AUTHORIZATION. Must run AFTER AuthRequired. Checks that the
// verified token carries the given permission.
// Failure → 403: we know who you are, but you may not do this.
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, _ := c.Get("scopes")
		scopes, _ := v.([]string)
		for _, s := range scopes {
			if s == scope {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "missing required scope: " + scope,
		})
	}
}

// CurrentSub returns the verified user id inside handlers (or "" if absent).
func CurrentSub(c *gin.Context) string {
	if v, ok := c.Get("sub"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}