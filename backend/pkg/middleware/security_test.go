package middleware

import (
	"github.com/fabiokaelin/f-oauth/config"
	"github.com/gin-gonic/gin"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIdentityWritesRejectOtherOrigins(t *testing.T) {
	old := config.FrontEndOrigin
	config.FrontEndOrigin = "https://oauth.fabkli.ch"
	defer func() { config.FrontEndOrigin = old }()
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		origin, cookie string
		status         int
	}{{"https://oauth.fabkli.ch", "token=opaque", 200}, {"https://tipp.fabkli.ch", "token=opaque", 403}, {"https://attacker.example", "token=opaque", 403}, {"", "token=opaque", 403}, {"", "", 200}} {
		r := gin.New()
		r.Use(ProtectWrites())
		r.POST("/", func(c *gin.Context) { c.Status(200) })
		req := httptest.NewRequest("POST", "https://oauth.fabkli.ch/", strings.NewReader("{}"))
		req.Header.Set("Origin", test.origin)
		req.Header.Set("Cookie", test.cookie)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != test.status {
			t.Fatalf("origin %q got %d", test.origin, w.Code)
		}
	}
}

func TestProfileImageLimitDoesNotUseSmallJSONLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		path         string
		size, status int
	}{
		{"/api/users/me/image", 20000, 200}, {"/api/users/me/image", 11 << 20, 413}, {"/api/auth/login", 20000, 413},
	} {
		r := gin.New()
		r.Use(ProtectWrites())
		r.POST(test.path, func(c *gin.Context) {
			if _, err := io.Copy(io.Discard, c.Request.Body); err != nil {
				c.Status(413)
			} else {
				c.Status(200)
			}
		})
		q := httptest.NewRequest("POST", test.path, strings.NewReader(strings.Repeat("x", test.size)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		if w.Code != test.status {
			t.Fatalf("%s: got %d", test.path, w.Code)
		}
	}
}
