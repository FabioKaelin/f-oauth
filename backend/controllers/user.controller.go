package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
	"github.com/wpcodevo/google-github-oath2-golang/models"
	"github.com/wpcodevo/google-github-oath2-golang/utils"
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
	fmt.Println("ifExist", ifExist)
	if !ifExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	// update user
	_, err = utils.RunSQL("UPDATE `users` SET `name` = ?, `photo` = ? WHERE `id` = ?", currentUser.Name, currentUser.Photo, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "update user in database failed"})
		return
	}

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
	filePath := "http://localhost:8001/images/" + newFileName

	imageFile, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imageFile); err != nil {
		log.Fatal(err)
	}

	pngFile, err := png.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	src := imaging.Fill(pngFile, 400, 400, imaging.Center, imaging.Lanczos)
	// src := imaging.Fill(pngFile, 100, 100, imaging.Center, imaging.Lanczos)
	// src := imaging.Resize(pngFile, 1000, 0, imaging.Lanczos)
	err = imaging.Save(src, fmt.Sprintf("public/images/%v", newFileName))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"filepath": filePath})
}

func GetProfileImage(ctx *gin.Context) {
	userID := ctx.Param("userid")
	fmt.Println("userID", userID)
	// check if public/images/profileimage-userID.png exists
	// if not, return default image (public/default.png)
	// if yes, return public/images/profileimage-userID.png

	if _, err := os.Stat("public/images/profileimage-" + userID + ".png"); err == nil {
		fmt.Println("file exists")
		ctx.Status(http.StatusOK)
		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.File("public/images/profileimage-" + userID + ".png")

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.Status(http.StatusOK)
		ctx.File("public/default.png")
	} else {
		fmt.Println("error durring checking if file exists")
		// ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error durring checking if file exists"})
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.Status(http.StatusOK)
		ctx.File("public/default.png")
	}

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
		fmt.Println(err)
	}

	// Make request with marshalled JSON as the POST body
	_, err = http.Post(url, "application/json",
		bytes.NewBuffer(userJSON))

	if err != nil {
		fmt.Println(err)
	}

}
