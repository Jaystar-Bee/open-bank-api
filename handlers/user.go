package handlers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/jwt"
	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/Jaystar-Bee/open-bank-api/utils"
	"github.com/gin-gonic/gin"
)

func sendOTP(name, email string) error {
	arrayOfNumbers := utils.GenerateUniqueNumbers(1, 99999)
	otp := utils.JoinIntSlice(arrayOfNumbers)
	date := time.Now().Format(time.RFC822)
	template_data := map[string]any{
		"OTP":      otp,
		"Date":     date,
		"Name":     name,
		"Help":     os.Getenv("EMAIL_ACCOUNT"),
		"HelpLink": "mailto:" + os.Getenv("EMAIL_ACCOUNT"),
	}

	body, err := utils.ParseTemplate("emails/templates/otp.html", template_data)
	if err != nil {
		return err
	}
	messageStatus := make(chan bool)
	go utils.SendEmail(email, "Verify your OTP", body, messageStatus)
	db.RDB.Set(db.Ctx, email, otp, time.Minute+10)
	if <-messageStatus {
		return nil
	} else {
		return errors.New("unable to send email")
	}
}

// CreateTags 		godoc
//
//	@Tags			User
//	@Description	Onboard user to the application.
//	@Summary		Create a user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.USER_REQUEST					true	"Create User"
//	@Success		201		{object}	models.HTTP_MESSAGE_ONLY_RESPONSE	"User created successfully"
//	@Failure		400		{object}	models.Error						"Check body"
//	@Failure		404		{object}	models.Error						"User not found"
//	@Failure		500		{object}	models.Error						"Internal server error"
//	@Router			/user/signup [post]
func CreateUser(context *gin.Context) {
	var user models.USER

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	if !utils.IsValidEmail(user.Email) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email",
		})
		return
	}
	if !utils.IsConvertibleToNumber(user.TransactionPin) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Transaction should be in digits",
		})
		return
	}

	_, err = models.GetUserByEmail(user.Email)
	if err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User already exists with this email",
		})
		return
	}
	if user.Phone != "" {
		_, err = models.GetUserByPhone(user.Phone)
		if err == nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "User already exists with your phone number",
			})
			return
		}
	}
	_, err = models.GetUserByTag(user.Tag)
	if err == nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User already exists with this tag",
		})
		return
	}
	err = user.Save()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to Save user",
			"dev_reason": err.Error(),
		})
		return
	}
	err = sendOTP(user.FirstName+" "+user.LastName, user.Email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to send OTP",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

// LogUserIn		godoc
//
//	@Tags			User
//	@Description	Log user in to the application.
//	@Summary		Log user in
//	@Accept			json
//	@Produce		json
//	@Param			login	body		models.USER_LOGIN			true	"Log User In"
//	@Success		200		{object}	models.HTTP_LOGIN_RESPONSE	"User logged in successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Router			/user/login [post]
func Login(context *gin.Context) {
	var login models.USER_LOGIN
	err := context.ShouldBindJSON(&login)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	// CHECK IF USER EXISTS WITH EMAIL
	user, err := models.GetUserByEmail(login.Email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Unable to find user, please check your email",
			"dev_reason": err.Error(),
		})
		return
	}

	// GENERATE TOKEN
	token, err := jwt.GenerateJWT(user.ID, user.Email, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to generate token",
			"dev_reason": err.Error(),
		})
		return
	}

	// LOG USER IN
	err = login.Login()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to login, Please check your password",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"data":    user,
		"token":   token,
	})

}

// RenewToken		godoc
//
//	@Tags			User
//	@Description	Renew user token.
//	@Summary		Renew token
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	models.HTTP_TOKEN_RESPONSE	"Token renewd successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/renew [get]
func RenewToken(context *gin.Context) {
	email := context.GetString("email")

	user, err := models.GetUserByEmail(email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Unable to find user",
		})
		return
	}

	token, err := jwt.GenerateJWT(user.ID, user.Email, time.Now().Add(time.Hour*6).Unix())

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to generate token",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Token renewed successfully",
		"token":   token,
	})

}

// GetUserByTag	godoc
//
//	@Tags			User
//	@Description	Get user by tag.
//	@Summary		Get user by tag.
//	@Accept			json
//	@Produce		json
//	@Param			tag	path		string						true	"User Tag"
//	@Success		200	{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/tag/{tag} [get]
func GetUserByTag(context *gin.Context) {
	tag := context.Param("tag")

	if tag == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User tag is required",
		})
		return
	}

	user, err := models.GetUserByTag(tag)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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

// GetUserByEmail	godoc
//
//	@Tags			User
//	@Description	Get user by email.
//	@Summary		Get user by email.
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string						true	"User Email"
//	@Success		200	{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/email/{email} [get]
func GetUserByEmail(context *gin.Context) {
	email := context.Param("email")

	if email == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User email is required",
		})
		return
	}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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

// GetUserByPhone	godoc
//
//	@Tags			User
//	@Description	Get user by phone.
//	@Summary		Get user by phone.
//	@Accept			json
//	@Produce		json
//	@Param			phone	path		string						true	"User Phone"
//	@Success		200	{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/phone/{phone} [get]
func GetUserByPhone(context *gin.Context) {
	phone := context.Param("phone")

	if phone == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User phone is required",
		})
		return
	}

	user, err := models.GetUserByPhone(phone)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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

// GetUserById	godoc
//
//	@Tags			User
//	@Description	Get user by Id.
//	@Summary		Get user by Id.
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string						true	"User Id"
//	@Success		200	{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/{user_id} [get]
func GetUserById(context *gin.Context) {
	id := context.Param("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to parse user id",
		})
		return
	}
	user, err := models.GetUserByID(parsedId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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
