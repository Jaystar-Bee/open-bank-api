package handlers

import (
	"net/http"

	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/gin-gonic/gin"
)

func CreateUser(context *gin.Context) {
	var user models.USER

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	_, err = models.GetUserByEmail(user.Email)
	if err == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}
