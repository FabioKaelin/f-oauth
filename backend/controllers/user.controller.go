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

	"github.com/fabiokaelin/f-oauth/initializers"
	"github.com/fabiokaelin/f-oauth/models"
	"github.com/fabiokaelin/f-oauth/utils"
	"github.com/gin-gonic/gin"
)

func GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	ctx.JSON(http.StatusOK, models.FilteredResponse(&currentUser))
	// ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(&currentUser)}})
}

func UpdateMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// read bodyUser from request
	var bodyUser models.User
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
	rows, err := utils.RunSQL("SELECT `id`, `name`, `email`, `password`, `role`, `photo`, `verified`, `provider`, `created_at`, `updated_at` FROM `users` WHERE `id` = ? LIMIT 1", currentUser.ID)
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
	rows, err = utils.RunSQL("UPDATE `users` SET `name` = ?, `photo` = ? WHERE `id` = ?", currentUser.Name, currentUser.Photo, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "update user in database failed"})
		return
	}
	rows.Close()

	updateInAllFProducts(currentUser)

	ctx.JSON(http.StatusOK, models.FilteredResponse(&currentUser))
}
func UploadResizeSingleFile(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

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
	req, err := http.NewRequest("POST", initializers.StartConfig.InternalImageService+"/api/users/"+currentUser.ID.String(), &b)
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

func GetProfileImage(ctx *gin.Context) {
	userID := ctx.Param("userid")
	fmt.Println("userID", userID)

	url := initializers.StartConfig.InternalImageService + "/api/users/" + userID
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "get image from image service failed"})
		return
	}
	defer resp.Body.Close()

	// Decode the image
	imageFile, _, err := image.Decode(resp.Body)
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

	ctx.Data(http.StatusOK, "image/png", buf.Bytes())

}

func updateInAllFProducts(currentUser models.User) {
	// Tipp
	url := initializers.StartConfig.InternalTippURL
	putRequest(url+"/internal/user", currentUser)

	// DevTipp
	url = initializers.StartConfig.InternalDevTippURL
	putRequest(url+"/internal/user", currentUser)
}

func putRequest(url string, currentUser models.User) {
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
