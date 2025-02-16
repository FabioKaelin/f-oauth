package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/middleware"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
)

func UserRouter(apiGroup *gin.RouterGroup) {
	userGroup := apiGroup.Group("/users")
	{
		userGroup.GET("/me", middleware.SetUserToContext(), userGetMe)
		userGroup.PUT("/me", middleware.SetUserToContext(), userPutMe)
		userGroup.POST("/me/image", middleware.SetUserToContext(), userPostMeImage)
		userGroup.GET("/:userid/image", userGetProfileImage)
	}
}

func userGetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(user_pkg.User)

	ctx.JSON(http.StatusOK, user_pkg.FilteredResponse(&currentUser))
	// ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": user_pkg.FilteredResponse(&currentUser)}})
}

func userPutMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(user_pkg.User)

	// read bodyUser from request
	var bodyUser user_pkg.User
	if err := ctx.ShouldBindJSON(&bodyUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if bodyUser.Name != "" {
		currentUser.Name = bodyUser.Name
	}

	if bodyUser.Photo != "" {
		currentUser.Photo = bodyUser.Photo
	}

	fmt.Println("id", currentUser.ID)
	rows, err := db.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `id` = ? LIMIT 1", currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message1": err.Error()})
		return
	}

	ifExist := rows.Next()
	rows.Close()
	fmt.Println("ifExist", ifExist)
	if !ifExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	// update user
	rows, err = db.RunSQL("UPDATE `users` SET `name` = ?, `photo` = ? WHERE `id` = ?", currentUser.Name, currentUser.Photo, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "update user in database failed"})
		return
	}
	rows.Close()

	updateInAllFProducts(currentUser)

	ctx.JSON(http.StatusOK, user_pkg.FilteredResponse(&currentUser))
}

func userPostMeImage(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(user_pkg.User)

	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	newFileName := "profileimage-" + currentUser.ID.String() + ".png"
	fmt.Println("newFileName", newFileName)

	imageFile, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "image decode failed"})
		return
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imageFile); err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "png encode failed"})
		return
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	defer w.Close()

	// Create a new form file
	fw, err := w.CreateFormFile("image", newFileName)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "create form file failed"})
		return
	}

	// Write the image data to the form file
	if _, err = io.Copy(fw, buf); err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "copy image data to form file failed"})
		return
	}

	// Close the multipart writer
	if err = w.Close(); err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "close multipart writer failed"})
		return
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", config.InternalImageService+"/api/users/"+currentUser.ID.String(), &b)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "create new http request failed"})
		return
	}

	// Set the content type, this is very important
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Do the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "do request failed"})
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("bad status: %s\n", res.Status)
	}

	ctx.JSON(http.StatusOK, gin.H{"worked": true})
}

func userGetProfileImage(ctx *gin.Context) {
	userID := ctx.Param("userid")
	fmt.Println("userID", userID)

	url := config.InternalImageService + "/api/profile/" + userID
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "get image from image service failed"})
		return
	}
	defer resp.Body.Close()

	ctx.DataFromReader(http.StatusOK, resp.ContentLength, "image/png", resp.Body, nil)

}

func updateInAllFProducts(currentUser user_pkg.User) {
	// Tipp
	url := config.InternalTippURL
	putRequest(url+"/internal/user", currentUser)

	// DevTipp
	url = config.InternalDevTippURL
	putRequest(url+"/internal/user", currentUser)
}

func putRequest(url string, currentUser user_pkg.User) {
	// Marshal it into JSON prior to requesting
	userJSON, err := json.Marshal(currentUser)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	// Make request with marshalled JSON as the POST body
	_, err = http.Post(url, "application/json",
		bytes.NewBuffer(userJSON))

	if err != nil {
		fmt.Println("error", err)
		return
	}

}
