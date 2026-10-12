package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOAuthStateAndRedirectBoundaries(t *testing.T) {
	for _, raw := range []string{"https://evil.example/", "https://fabkli.ch.evil.example/", "https://attacker@tipp.fabkli.ch/", "http://tipp.fabkli.ch/", "//tipp.fabkli.ch/", "https://tipp.fabkli.ch:999/"} {
		if allowedReturn(raw) {
			t.Fatalf("unsafe return URL accepted: %s", raw)
		}
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/start/:provider", startOAuth)
	router.GET("/callback", func(c *gin.Context) {
		from, ok := oauthReturn(c, "github")
		if !ok {
			c.AbortWithStatus(400)
			return
		}
		c.String(200, from)
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "https://oauth.fabkli.ch/start/github?from="+url.QueryEscape("https://tipp.fabkli.ch/"), nil))
	if w.Code != 303 {
		t.Fatalf("start status %d", w.Code)
	}
	location, _ := url.Parse(w.Header().Get("Location"))
	state := location.Query().Get("state")
	if len(state) != 43 {
		t.Fatal("OAuth state is not random 256-bit nonce")
	}
	req := httptest.NewRequest("GET", "https://oauth.fabkli.ch/callback?state="+state, nil)
	for _, cookie := range w.Result().Cookies() {
		if cookie.Domain != "" || !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
			t.Error("weak OAuth state cookie")
		}
		req.AddCookie(cookie)
	}
	result := httptest.NewRecorder()
	router.ServeHTTP(result, req)
	if result.Code != 200 || result.Body.String() != "https://tipp.fabkli.ch/" {
		t.Fatal("valid callback failed")
	}
	missing := httptest.NewRecorder()
	router.ServeHTTP(missing, httptest.NewRequest("GET", "https://oauth.fabkli.ch/callback?state="+state, nil))
	if missing.Code != 400 {
		t.Fatal("callback without browser-bound state accepted")
	}
	forged := httptest.NewRecorder()
	req.URL.RawQuery = "state=forged"
	router.ServeHTTP(forged, req)
	if forged.Code != 400 {
		t.Fatal("forged state accepted")
	}
}
