package middleware

import (
	"errors"
	"fmt"
	"github.com/fabiokaelin/fcommon/pkg/ratelimit"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/session"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
)

func SetUserToContext(renew ...bool) gin.HandlerFunc {
	doRenew := len(renew) == 0 || renew[0]
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-store")
		token := ""
		fields := strings.Fields(ctx.GetHeader("Authorization"))
		if len(fields) != 0 {
			if len(fields) != 2 || fields[0] != "Bearer" {
				ctx.AbortWithStatus(401)
				return
			}
			token = fields[1]
		} else {
			cookies := ctx.Request.CookiesNamed("token")
			if len(cookies) != 1 {
				ctx.AbortWithStatus(401)
				return
			}
			token = cookies[0].Value
		}
		validate := session.Validate
		if !doRenew {
			validate = session.Check
		}
		id, err := validate(token)
		if err != nil {
			status := http.StatusServiceUnavailable
			if errors.Is(err, db.ErrInvalidSession) {
				status = http.StatusUnauthorized
			}
			ctx.AbortWithStatusJSON(status, gin.H{"status": "fail", "message": "Authentication failed"})
			return
		}
		user, err := db.GetUserByID(id)
		if err != nil {
			ctx.AbortWithStatus(503)
			return
		}
		if user.ID == "" {
			ctx.AbortWithStatus(401)
			return
		}
		if !ratelimit.User(ctx, id) {
			return
		}
		if doRenew {
			session.SetCookie(ctx.Writer, token)
		}
		ctx.Set("sessionToken", token)
		ctx.Set("currentUser", user_pkg.ConvertDBUserToUser(user))
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Content-Type", "application/json")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			u, err := url.Parse(origin)
			if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" || (u.Hostname() != "fabkli.ch" && !strings.HasSuffix(u.Hostname(), ".fabkli.ch")) {
				c.AbortWithStatus(403)
				return
			}
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, Cookie, caches, Pragma, Expires")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		// Vary: Origin
		c.Writer.Header().Set("Vary", "Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func IPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/internal") {
			c.Next()
			return
		}

		ip := c.GetString("ratelimit.ip")
		if ip == "" {
			ip = c.ClientIP()
		}
		userAgent := c.Request.UserAgent()
		output := fmt.Sprintf("IP: '%s' - UserAgent: '%s'", ip, userAgent)
		fmt.Println(output)
		c.Next()
	}
}

// Browser writes to the identity service must originate from its own frontend.
// Non-browser bearer clients remain supported; body limits bound JSON processing.
func ProtectWrites() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" && c.Request.Method != "OPTIONS" {
			limit := int64(16384)
			if c.Request.Method == "POST" && c.Request.URL.Path == "/api/users/me/image" {
				limit = 10 << 20 // preserve the authenticated profile-image upload contract
			}
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
			origin := c.GetHeader("Origin")
			if origin != "" {
				if origin != config.FrontEndOrigin {
					c.AbortWithStatus(403)
					return
				}
			} else if len(c.Request.CookiesNamed("token")) > 0 {
				c.AbortWithStatus(403)
				return
			}
		}
		c.Next()
	}
}
