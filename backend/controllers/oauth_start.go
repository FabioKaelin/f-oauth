package controllers

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/session"
	"github.com/gin-gonic/gin"
)

func allowedReturn(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.User == nil && u.Port() == "" && (u.Hostname() == "fabkli.ch" || strings.HasSuffix(u.Hostname(), ".fabkli.ch")) && !strings.ContainsAny(raw, "\r\n\\")
}

func startOAuth(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider != "github" && provider != "google" {
		ctx.AbortWithStatus(404)
		return
	}
	from := ctx.Query("from")
	if from == "" {
		from = config.FrontEndOrigin + "/profile"
	}
	if !allowedReturn(from) || len(from) > 2048 {
		ctx.AbortWithStatus(400)
		return
	}
	nonce, err := session.NewToken()
	if err != nil {
		ctx.AbortWithStatus(503)
		return
	}
	cookie := func(name, value string) {
		http.SetCookie(ctx.Writer, &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: 600, Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	}
	// __Host- cookies cannot be set by other FabKli subdomains.
	cookie("__Host-f-oauth-state", provider+"."+nonce)
	cookie("__Host-f-oauth-return", base64.RawURLEncoding.EncodeToString([]byte(from)))
	q := url.Values{"state": {nonce}, "response_type": {"code"}}
	target := "https://github.com/login/oauth/authorize"
	if provider == "github" {
		q.Set("client_id", config.GitHubClientID)
		q.Set("redirect_uri", config.GitHubOAuthRedirectUrl)
		q.Set("scope", "user:email")
	} else {
		target = "https://accounts.google.com/o/oauth2/v2/auth"
		q.Set("client_id", config.GoogleClientID)
		q.Set("redirect_uri", config.GoogleOAuthRedirectUrl)
		q.Set("scope", "openid email profile")
	}
	ctx.Header("Cache-Control", "no-store")
	ctx.Redirect(http.StatusSeeOther, target+"?"+q.Encode())
}

func oauthReturn(ctx *gin.Context, provider string) (string, bool) {
	states := ctx.Request.CookiesNamed("__Host-f-oauth-state")
	returns := ctx.Request.CookiesNamed("__Host-f-oauth-return")
	for _, name := range []string{"__Host-f-oauth-state", "__Host-f-oauth-return"} {
		http.SetCookie(ctx.Writer, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1, Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	}
	state := ctx.Query("state")
	if len(states) != 1 || len(returns) != 1 || len(state) != 43 || subtle.ConstantTimeCompare([]byte(states[0].Value), []byte(provider+"."+state)) != 1 {
		return "", false
	}
	raw, err := base64.RawURLEncoding.DecodeString(returns[0].Value)
	if err != nil || !allowedReturn(string(raw)) {
		return "", false
	}
	return string(raw), true
}
