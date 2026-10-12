package controllers

import (
	"fmt"
	"net/http"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/github"
	"github.com/fabiokaelin/f-oauth/pkg/google"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	"github.com/fabiokaelin/f-oauth/pkg/session"

	"github.com/gin-gonic/gin"
)

func OAuth2Router(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/sessions/check", middleware.SetUserToContext(false), userGetMe)
	oauth2Group := apiGroup.Group("/sessions/oauth")
	{
		oauth2Group.GET("/google", oauth2Google)
		oauth2Group.GET("/github", oauth2GitHub)
		oauth2Group.GET("/start/:provider", startOAuth)
	}
}

func oauth2Google(ctx *gin.Context) {

	// TODO: Redirect to loginpage when an error occurs with error message
	code := ctx.Query("code")
	state, valid := oauthReturn(ctx, "google")
	if !valid {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "Invalid or expired login state; please start login again"})
		return
	}

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

	session.SetCookie(ctx.Writer, token)
	ctx.Redirect(http.StatusSeeOther, state)
}

func oauth2GitHub(ctx *gin.Context) {

	// TODO: Redirect to loginpage when an error occurs with error message
	code := ctx.Query("code")
	state, valid := oauthReturn(ctx, "github")
	if !valid {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "Invalid or expired login state; please start login again"})
		return
	}

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

	session.SetCookie(ctx.Writer, token)
	ctx.Redirect(http.StatusSeeOther, state)
}
