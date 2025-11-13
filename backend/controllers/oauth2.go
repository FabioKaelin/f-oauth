package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/github"
	"github.com/fabiokaelin/f-oauth/pkg/google"

	"github.com/gin-gonic/gin"
)

func OAuth2Router(apiGroup *gin.RouterGroup) {
	oauth2Group := apiGroup.Group("/sessions/oauth")
	{
		oauth2Group.GET("/google", oauth2Google)
		oauth2Group.GET("/github", oauth2GitHub)
	}
}

// isValidRedirectURL validates that the redirect URL is safe
func isValidRedirectURL(redirectURL string, allowedOrigin string) bool {
	if redirectURL == "" {
		return false
	}

	// Parse the redirect URL
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}

	// Parse the allowed origin
	allowedURL, err := url.Parse(allowedOrigin)
	if err != nil {
		return false
	}

	// Allow relative URLs (they start with /)
	if parsedURL.Scheme == "" && parsedURL.Host == "" {
		return strings.HasPrefix(redirectURL, "/")
	}

	// For absolute URLs, check if the host matches the allowed origin
	return parsedURL.Scheme == allowedURL.Scheme && parsedURL.Host == allowedURL.Host
}

func oauth2Google(ctx *gin.Context) {
	fmt.Println(ctx.Request.URL.String())
	// TODO: Redirect to loginpage when an error occurs with error message
	code := ctx.Query("code")
	state := ctx.Query("state") // path/url to redirect after login

	if code == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login"))
		return
	}

	token, errorcode, err := google.LoginWithCode(code)

	if err != nil {
		fmt.Println("error, failed login with google:", err)
		if errorcode == 2 {
			ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=already_signed_up_with_different_method"))
			return
			// ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		} else {
			ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=error_occured"))
			// ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "something went wrong"})
			return
		}
	}

	if strings.Contains(config.TokenURL, "localhost") {
		ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", true, true)
	} else {
		ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", config.TokenURL, true, true)
	}
	
	// Validate and sanitize the redirect URL
	redirectUrl := ""
	if state != "" && isValidRedirectURL(state, config.FrontEndOrigin) {
		// If state is a relative path, prepend the frontend origin
		if strings.HasPrefix(state, "/") {
			redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, state)
		} else {
			// If it's an absolute URL that passed validation, use it
			redirectUrl = state
		}
	} else {
		// Default to profile page if state is invalid or empty
		redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, "/profile")
	}
	fmt.Println("redirectUrl", redirectUrl)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func oauth2GitHub(ctx *gin.Context) {
	fmt.Println(ctx.Request.URL.String())
	// TODO: Redirect to loginpage when an error occurs with error message
	code := ctx.Query("code")
	state := ctx.Query("state") // path/url to redirect after login

	if code == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login"))
		return
	}

	token, errorcode, err := github.LoginWithCode(code)

	if err != nil {
		fmt.Println("error, failed login with google:", err)
		if errorcode == 2 {
			ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=already_signed_up_with_different_method"))
			return
			// ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have already signed up with a different method"})
		} else {
			ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", config.FrontEndOrigin, "/login?error=error_occured"))
			// ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "something went wrong"})
			return
		}
	}

	if strings.Contains(config.TokenURL, "localhost") {
		ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", true, true)
	} else {
		ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", config.TokenURL, true, true)
	}
	
	// Validate and sanitize the redirect URL
	redirectUrl := ""
	if state != "" && isValidRedirectURL(state, config.FrontEndOrigin) {
		// If state is a relative path, prepend the frontend origin
		if strings.HasPrefix(state, "/") {
			redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, state)
		} else {
			// If it's an absolute URL that passed validation, use it
			redirectUrl = state
		}
	} else {
		// Default to profile page if state is invalid or empty
		redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, "/profile")
	}
	fmt.Println("redirectUrl", redirectUrl)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
