package controllers

import (
	userpkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProfilePasswordRejectsMissingCurrentOrExternalProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct{ provider, body string }{{"local", `{"password":"replacement"}`}, {"google", `{"password":"replacement","current_password":"old"}`}, {"local", `null`}} {
		r := gin.New()
		r.POST("/change", func(c *gin.Context) { c.Set("currentUser", userpkg.User{Provider: test.provider}); changePassword(c) })
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/change", strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Fatalf("provider %s returned %d", test.provider, w.Code)
		}
	}
}
