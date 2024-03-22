package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/Jaystar-Bee/open-bank-api/utils"
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
	if !utils.IsValidEmail(user.Email) {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email",
		})
		return
	}
	if !utils.IsConvertibleToNumber(user.TransactionPin) {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Transaction should be in digits",
		})
		return
	}

	_, err = models.GetUserByEmail(user.Email)
	if err == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists with this email",
		})
		return
	}
	if user.Phone != "" {
		_, err = models.GetUserByPhone(user.Phone)
		if err == nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "User already exists with your phone number",
			})
			return
		}
	}
	_, err = models.GetUserByTag(user.Tag)
	if err == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists with this tag",
		})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to Save user",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func GetUserByTag(context *gin.Context) {
	tag := context.Param("tag")

	if tag == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User tag is required",
		})
		return
	}

	user, err := models.GetUserByTag(tag)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message":    "Unable to find user",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully",
		"data":    user,
	})
}

func GetUserByEmail(context *gin.Context) {
	email := context.Param("email")

	if email == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User email is required",
		})
		return
	}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message":    "Unable to find user",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully",
		"data":    user,
	})
}

func GetUserByPhone(context *gin.Context) {
	phone := context.Param("phone")

	if phone == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "User phone is required",
		})
		return
	}

	user, err := models.GetUserByPhone(phone)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message":    "Unable to find user",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully",
		"data":    user,
	})
}
func GetUserById(context *gin.Context) {
	id := context.Param("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to parse user id",
		})
		return
	}
	user, err := models.GetUserByID(parsedId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message":    "Unable to find user",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully",
		"data":    user,
	})
}
