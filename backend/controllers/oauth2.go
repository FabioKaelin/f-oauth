package controllers

import (
	"fmt"
	"net/http"
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
	redirectUrl := ""
	// if pathUrl not begin with http or https
	fmt.Println("pathUrl", state)
	if !strings.HasPrefix(state, "http") && !strings.HasPrefix(state, "https") {
		state = "/profile"
		redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, state)
	} else {
		redirectUrl = state
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
	redirectUrl := ""
	// if pathUrl not begin with http or https
	fmt.Println("pathUrl", state)
	if !strings.HasPrefix(state, "http") && !strings.HasPrefix(state, "https") {
		state = "/profile"
		redirectUrl = fmt.Sprintf("%s%s", config.FrontEndOrigin, state)
	} else {
		redirectUrl = state
	}
	fmt.Println("redirectUrl", redirectUrl)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
