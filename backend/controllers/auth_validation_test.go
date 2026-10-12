package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidAuthPayloadsDoNotReachDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct{ path, body string }{{"/register", "null"}, {"/login", "null"}, {"/register", `{"name":"Valid User","email":"test@example.com","password":"x"}`}, {"/register", `{"name":"Valid User","email":"test@example.com","password":"` + strings.Repeat("a", 73) + `"}`}, {"/login", `{"password":"valid-password"}`}} {
		r := gin.New()
		AuthRouter(r.Group("/api/auth-test"))
		path := "/api/auth-test/auth" + test.path
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", path, strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Fatalf("%s body %s returned %d", test.path, test.body, w.Code)
		}
	}
}
