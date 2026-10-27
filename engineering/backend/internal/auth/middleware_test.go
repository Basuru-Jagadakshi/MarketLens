package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthRequired_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", AuthRequired(), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthRequired_NonBearerHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", AuthRequired(), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	resetJWKSCache(t)
	withThunderConfig(t, "http://unused.invalid", "issuer", "aud")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", AuthRequired(), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthRequired_ValidToken_SetsContext(t *testing.T) {
	resetJWKSCache(t)
	key := generateTestKey(t)
	srv := newJWKSServer(t, []jwk{jwkFromPublicKey("kid1", &key.PublicKey)})
	withThunderConfig(t, srv.URL, "issuer", "aud")

	raw := signTestToken(t, key, "kid1", validTestClaims("issuer", "aud", "user-42", "book:create"))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	var gotSub string
	var gotScopes []string
	r.GET("/x", AuthRequired(), func(c *gin.Context) {
		gotSub = CurrentSub(c)
		v, _ := c.Get("scopes")
		gotScopes, _ = v.([]string)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
	if gotSub != "user-42" {
		t.Errorf("expected sub %q, got %q", "user-42", gotSub)
	}
	if len(gotScopes) != 1 || gotScopes[0] != "book:create" {
		t.Errorf("unexpected scopes: %v", gotScopes)
	}
}

func TestRequireScope_HasScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x",
		func(c *gin.Context) { c.Set("scopes", []string{"book:create", "book:update"}); c.Next() },
		RequireScope("book:update"),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRequireScope_MissingScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x",
		func(c *gin.Context) { c.Set("scopes", []string{"book:create"}); c.Next() },
		RequireScope("book:delete"),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestRequireScope_NoScopesSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", RequireScope("book:create"), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestCurrentSub(t *testing.T) {
	gin.SetMode(gin.TestMode)

	empty, _ := gin.CreateTestContext(httptest.NewRecorder())
	if got := CurrentSub(empty); got != "" {
		t.Errorf("expected empty string when sub is not set, got %q", got)
	}

	withSub, _ := gin.CreateTestContext(httptest.NewRecorder())
	withSub.Set("sub", "user-1")
	if got := CurrentSub(withSub); got != "user-1" {
		t.Errorf("expected %q, got %q", "user-1", got)
	}
}
