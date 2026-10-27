package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// test helpers shared with middleware_test.go

func generateTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating test RSA key: %v", err)
	}
	return key
}

func jwkFromPublicKey(kid string, pub *rsa.PublicKey) jwk {
	eBytes := big.NewInt(int64(pub.E)).Bytes()
	return jwk{
		Kid: kid,
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(eBytes),
	}
}

func newJWKSServer(t *testing.T, keys []jwk) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": keys})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func signTestToken(t *testing.T, key *rsa.PrivateKey, kid string, claims customClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("signing test token: %v", err)
	}
	return signed
}

// withThunderConfig overrides the package-level ThunderID config for the
// duration of the test, restoring the original values afterward.
func withThunderConfig(t *testing.T, baseURL, issuer, audience string) {
	t.Helper()
	origBase, origIssuer, origAudience := thunderBaseURL, thunderIssuer, thunderAudience
	thunderBaseURL, thunderIssuer, thunderAudience = baseURL, issuer, audience
	t.Cleanup(func() {
		thunderBaseURL, thunderIssuer, thunderAudience = origBase, origIssuer, origAudience
	})
}

// resetJWKSCache clears the package-level JWKS cache before and after the
// test, since it's global singleton state shared across every test in
// this package.
func resetJWKSCache(t *testing.T) {
	t.Helper()
	clear := func() {
		jwksMu.Lock()
		jwksCache = nil
		jwksLastRefresh = time.Time{}
		jwksMu.Unlock()
	}
	clear()
	t.Cleanup(clear)
}

// ValidateConfig 

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name                      string
		baseURL, issuer, audience string
		wantErrContains           string
	}{
		{"all set", "http://thunder", "http://thunder", "aud", ""},
		{"missing base url", "", "http://thunder", "aud", "THUNDER_BASE_URL"},
		{"missing issuer", "http://thunder", "", "aud", "THUNDER_ISSUER"},
		{"missing audience", "http://thunder", "http://thunder", "", "THUNDER_AUDIENCE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withThunderConfig(t, tt.baseURL, tt.issuer, tt.audience)

			err := ValidateConfig()
			if tt.wantErrContains == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErrContains, err)
			}
		})
	}
}

// translatePublicKeyFromJWK

func TestTranslatePublicKeyFromJWK_RoundTrips(t *testing.T) {
	key := generateTestKey(t)
	j := jwkFromPublicKey("kid1", &key.PublicKey)

	pub, err := translatePublicKeyFromJWK(j)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pub.E != key.PublicKey.E || pub.N.Cmp(key.PublicKey.N) != 0 {
		t.Errorf("translated key does not match the original")
	}
}

func TestTranslatePublicKeyFromJWK_BadBase64(t *testing.T) {
	_, err := translatePublicKeyFromJWK(jwk{N: "not valid base64!!", E: "AQAB"})
	if err == nil {
		t.Fatal("expected an error for an undecodable N")
	}
}

func TestTranslatePublicKeyFromJWK_OversizedExponent(t *testing.T) {
	eBytes := make([]byte, 9) // more than the 8 bytes this function supports
	for i := range eBytes {
		eBytes[i] = 1
	}
	j := jwk{
		N: base64.RawURLEncoding.EncodeToString([]byte{1, 2, 3}),
		E: base64.RawURLEncoding.EncodeToString(eBytes),
	}

	_, err := translatePublicKeyFromJWK(j)
	if err == nil {
		t.Fatal("expected an error for an oversized exponent")
	}
}

// findJWK 

func TestFindJWK(t *testing.T) {
	keys := []jwk{{Kid: "a"}, {Kid: "b"}}

	if k, ok := findJWK(keys, "b"); !ok || k.Kid != "b" {
		t.Errorf("expected to find kid %q", "b")
	}
	if _, ok := findJWK(keys, "missing"); ok {
		t.Error("expected not to find an unknown kid")
	}
}

// fetchJWKSLocked / getJWKS / refreshJWKS

func TestGetJWKS_Success(t *testing.T) {
	resetJWKSCache(t)

	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "issuer", "aud")

	keys, err := getJWKS()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 || keys[0].Kid != "kid1" {
		t.Errorf("unexpected keys: %+v", keys)
	}
}

func TestFetchJWKSLocked_NonOKStatus_DoesNotUpdateCache(t *testing.T) {
	resetJWKSCache(t)
	jwksMu.Lock()
	jwksCache = []jwk{{Kid: "existing"}}
	jwksMu.Unlock()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	withThunderConfig(t, srv.URL, "issuer", "aud")

	jwksMu.Lock()
	_, err := fetchJWKSLocked()
	jwksMu.Unlock()
	if err == nil {
		t.Fatal("expected an error for a non-200 response")
	}

	jwksMu.Lock()
	defer jwksMu.Unlock()
	if len(jwksCache) != 1 || jwksCache[0].Kid != "existing" {
		t.Errorf("cache should be untouched after a failed fetch, got %+v", jwksCache)
	}
}

func TestFetchJWKSLocked_EmptyKeys_Rejected(t *testing.T) {
	resetJWKSCache(t)
	jwksMu.Lock()
	jwksCache = []jwk{{Kid: "existing"}}
	jwksMu.Unlock()

	srv := newJWKSServer(t, []jwk{})
	withThunderConfig(t, srv.URL, "issuer", "aud")

	jwksMu.Lock()
	_, err := fetchJWKSLocked()
	jwksMu.Unlock()
	if err == nil {
		t.Fatal("expected an error for an empty keys list")
	}

	jwksMu.Lock()
	defer jwksMu.Unlock()
	if len(jwksCache) != 1 || jwksCache[0].Kid != "existing" {
		t.Errorf("cache should be untouched after a rejected fetch, got %+v", jwksCache)
	}
}

func TestRefreshJWKS_RespectsCooldown(t *testing.T) {
	resetJWKSCache(t)

	var hits int
	key := generateTestKey(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []jwk{jwkFromPublicKey("kid1", &key.PublicKey)}})
	}))
	t.Cleanup(srv.Close)
	withThunderConfig(t, srv.URL, "issuer", "aud")

	if _, err := refreshJWKS(); err != nil {
		t.Fatalf("unexpected error on first refresh: %v", err)
	}
	if _, err := refreshJWKS(); err != nil {
		t.Fatalf("unexpected error on second refresh: %v", err)
	}
	if hits != 1 {
		t.Errorf("expected the cooldown to prevent a second fetch, got %d hits", hits)
	}
}

// verifyToken (full signature + claims verification)

func validTestClaims(issuer, audience, subject, scope string) customClaims {
	now := time.Now()
	return customClaims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
}

func TestVerifyToken_Valid(t *testing.T) {
	resetJWKSCache(t)
	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "test-issuer", "test-audience")

	raw := signTestToken(t, key, "kid1", validTestClaims("test-issuer", "test-audience", "user-1", "book:create book:update"))

	claims, err := verifyToken(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Errorf("expected subject %q, got %q", "user-1", claims.Subject)
	}
	if claims.Scope != "book:create book:update" {
		t.Errorf("unexpected scope: %q", claims.Scope)
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	resetJWKSCache(t)
	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "test-issuer", "test-audience")

	claims := validTestClaims("test-issuer", "test-audience", "user-1", "")
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
	raw := signTestToken(t, key, "kid1", claims)

	if _, err := verifyToken(raw); err == nil {
		t.Fatal("expected an error for an expired token")
	}
}

func TestVerifyToken_WrongIssuer(t *testing.T) {
	resetJWKSCache(t)
	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "test-issuer", "test-audience")

	raw := signTestToken(t, key, "kid1", validTestClaims("someone-else", "test-audience", "user-1", ""))

	if _, err := verifyToken(raw); err == nil {
		t.Fatal("expected an error for a mismatched issuer")
	}
}

func TestVerifyToken_WrongAudience(t *testing.T) {
	resetJWKSCache(t)
	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "test-issuer", "test-audience")

	raw := signTestToken(t, key, "kid1", validTestClaims("test-issuer", "someone-else", "user-1", ""))

	if _, err := verifyToken(raw); err == nil {
		t.Fatal("expected an error for a mismatched audience")
	}
}

func TestVerifyToken_UnknownKid(t *testing.T) {
	resetJWKSCache(t)
	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "test-issuer", "test-audience")

	raw := signTestToken(t, key, "kid-does-not-exist", validTestClaims("test-issuer", "test-audience", "user-1", ""))

	if _, err := verifyToken(raw); err == nil {
		t.Fatal("expected an error for an unrecognized kid")
	}
}

func TestVerifyToken_BadSignature(t *testing.T) {
	resetJWKSCache(t)
	realKey := generateTestKey(t)
	wrongKey := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &realKey.PublicKey)})
	withThunderConfig(t, srv.URL, "test-issuer", "test-audience")

	// Signed with wrongKey but claims kid1, whose published JWK is realKey's
	// public key — the signature must not verify.
	raw := signTestToken(t, wrongKey, "kid1", validTestClaims("test-issuer", "test-audience", "user-1", ""))

	if _, err := verifyToken(raw); err == nil {
		t.Fatal("expected an error for a signature that doesn't match the published key")
	}
}
