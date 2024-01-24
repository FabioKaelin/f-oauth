package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
