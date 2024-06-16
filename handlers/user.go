package handlers

import (
	"errors"
	"fmt"
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

func sendOTP(name, email string) (int, error) {
	otp := utils.GenerateUniqueNumbers(1000, 9999)
	date := time.Now().Format(time.RFC822)
	expTime := 10
	template_data := map[string]any{
		"OTP":      otp,
		"Date":     date,
		"Time":     expTime,
		"Name":     name,
		"Help":     os.Getenv("EMAIL_ACCOUNT"),
		"HelpLink": "mailto:" + os.Getenv("EMAIL_ACCOUNT"),
	}

	body, err := utils.ParseTemplate("emails/templates/otp.html", template_data)
	if err != nil {
		return otp, err
	}
	messageStatus := make(chan bool)
	go utils.SendEmail(email, "Verify your Account", body, name, []string{"onboarding", "verify"}, messageStatus)
	db.RDB.Set(db.Ctx, email, otp, time.Minute+time.Duration(expTime))
	if <-messageStatus {
		return otp, nil
	} else {
		return otp, errors.New("unable to send email")
	}
}

// ToggleAccountDeactivation godoc
//
//	@Summary		Toggle account activation
//
//	@Description	Toggle account activation
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			otp	body		models.OTP							true	"OTP"
//	@Success		200	{object}	models.HTTP_MESSAGE_ONLY_RESPONSE	"Account deactivated successfully"
//	@Success		200	{object}	models.HTTP_LOGIN_RESPONSE			"Account activated successfully"
//	@Failure		400	{object}	models.Error						"Unable to process request"
//	@Failure		500	{object}	models.Error						"Unable to process request"
//	@Router			/user/toggle-account-deactivation [post]
func ToogleAccountDeactivation(context *gin.Context) {
	userId := context.GetInt64("user")
	var otpDetail models.OTP
	err := context.ShouldBindJSON(&otpDetail)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}

	// CHECK OTP
	otp_db, err := db.RDB.Get(db.Ctx, user.Email).Result()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "OTP Expired",
			"dev_reason": err.Error(),
		})
		return
	}

	if otpDetail.OTP != otp_db {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "OTP is invalid",
		})
		return
	}

	user.AccountIsDeactivated = !user.AccountIsDeactivated
	err = user.UpdateUser()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	if user.AccountIsDeactivated {
		context.JSON(http.StatusOK, gin.H{
			"message": "Account deactivated successfully",
		})
	} else {
		var token, err = jwt.GenerateJWT(user.ID, user.Email, time.Now().Add(time.Hour*24).Unix())
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message":    "Unable to process request",
				"dev_reason": err.Error(),
			})
			return
		}
		context.JSON(http.StatusAccepted, gin.H{
			"message": "Account activated successfully",
			"data":    user,
			"token":   token,
		})
	}
}

// VerifyAccount godoc
//
//	@Summary		Verify Account
//	@Description	Verify Account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			otp	body		models.OTP					true	"OTP"
//	@Success		200	{object}	models.HTTP_LOGIN_RESPONSE	"Account verified successfully"
//	@Failure		400	{object}	models.Error				"Unable to process request"
//	@Failure		500	{object}	models.Error				"Unable to process request"
//	@Router			/user/verify [post]
func VerifyAccount(context *gin.Context) {
	var otp models.OTP

	// CHECK OTP
	err := context.ShouldBindJSON(&otp)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	// CHECK OTP EXPIRATION AND IF USER EXIST
	email := otp.Email
	user, err := models.GetUserByEmail(email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}
	otp_db, err := db.RDB.Get(db.Ctx, email).Result()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "OTP expired",
			"dev_reason": err.Error(),
		})
		return
	}

	// CHECK OTP
	if otp_db != otp.OTP {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "OTP is invalid",
		})
		return
	}

	// VERIFY ACCOUNT AND UPDATE ACCOUNT
	user.IsVerified = true
	err = user.UpdateUser()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	// GENERATE JWT
	token, err := jwt.GenerateJWT(user.ID, user.Email, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	// SEND DATA
	context.JSON(http.StatusOK, gin.H{
		"message": "Account verified successfully",
		"data":    user,
		"token":   token,
	})
}

// SendOTP godoc
//
//	@summary		Send OTP
//	@description	Send OTP
//	@Tags			User
//	@accept			json
//	@produce		json
//	@param			user	body		models.OTP_REQUEST					true	"User"
//	@success		200		{object}	models.HTTP_MESSAGE_ONLY_RESPONSE	"OTP sent successfully"
//	@failure		400		{object}	models.Error						"Unable to process request"
//	@failure		400		{object}	models.Error						"Invalid email"
//	@failure		500		{object}	models.Error						"Unable to process request"
//	@router			/user/sendotp [post]
func SendOTP(context *gin.Context) {
	var userOTP models.OTP_REQUEST
	err := context.ShouldBindJSON(&userOTP)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	if !utils.IsValidEmail(userOTP.Email) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email",
		})
		return
	}
	user, _ := models.GetUserByEmail(userOTP.Email)
	var otp int
	if user.ID > 0 {
		otp, _ = sendOTP(fmt.Sprint(user.FirstName, " ", user.LastName), userOTP.Email)
	} else {
		otp, _ = sendOTP(userOTP.Name, userOTP.Email)
	}
	// if err != nil {
	// 	context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"message":    "Unable to send OTP",
	// 		"dev_reason": err.Error(),
	// 	})
	// 	return
	// }
	context.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"data":    otp,
	})
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
	otp, _ := sendOTP(user.FirstName+" "+user.LastName, user.Email)
	// if err != nil {
	// 	context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"message":    "Unable to send OTP",
	// 		"dev_reason": err.Error(),
	// 	})
	// 	return
	// }

	context.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"data":    otp,
	})
}

// EditUser godoc
//
//	@Summary		Edit User
//	@Description	Edit User
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.USER_EDIT					true	"Edit User"
//	@Success		200		{object}	models.HTTP_MESSAGE_ONLY_RESPONSE	"User updated successfully"
//	@Failure		400		{object}	models.Error						"Unable to process request"
//	@Failure		404		{object}	models.Error						"User not found"
//	@Failure		500		{object}	models.Error						"Unable to process request"
//	@Router			/user/edit [put]
func EditUser(context *gin.Context) {
	var editData models.USER_EDIT
	userId := context.GetInt64("user")

	err := context.ShouldBindJSON(&editData)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}
	user.FirstName = editData.FirstName
	user.LastName = editData.LastName
	user.Phone = editData.Phone
	user.Tag = editData.Tag
	user.Avatar = editData.Avatar

	err = user.UpdateUser()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
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
//	@Success		200		{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
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
//	@Success		200		{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
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
//	@Success		200		{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
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

// ChangeUserPin	godoc
//
//	@Tags			User
//	@Description	Change user pin.
//	@Summary		Change user pin.
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			pin	body		models.CHANGE_PIN			true	"User Pin"
//	@Success		200	{object}	models.HTTP_USER_RESPONSE	"User fetched successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/change-pin [patch]
func ChangeUserPin(context *gin.Context) {
	var pinData models.CHANGE_PIN
	userId := context.GetInt64("user")
	err := context.ShouldBindJSON(&pinData)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}

	err = user.ConfirmPin(pinData.OldPin)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid old pin",
			"dev_reason": err.Error(),
		})
		return
	}
	newPin, err := utils.HashText(pinData.NewPin)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to secure new pin",
			"dev_reason": err.Error(),
		})
		return
	}

	err = user.UpdatePin(newPin)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to update pin",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Pin updated successfully",
	})

}

// ChangeUserPassword godoc
//
//	@Summary		Change user password.
//	@Description	Change user password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			password	body		models.CHANGE_PASSWORD		true	"User Password"
//	@Success		200			{object}	models.HTTP_USER_RESPONSE	"Password updated successfully"
//	@Failure		400			{object}	models.Error
//	@Failure		404			{object}	models.Error
//	@Failure		500			{object}	models.Error
//	@Router			/user/change-password [patch]
func ChangeUserPassword(context *gin.Context) {
	var passwordData models.CHANGE_PASSWORD
	userId := context.GetInt64("user")
	err := context.ShouldBindJSON(&passwordData)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}
	err = user.ConfirmPin(passwordData.OldPassword)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid old password",
			"dev_reason": err.Error(),
		})
		return
	}

	newPassword, err := utils.HashText(passwordData.NewPassword)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to secure new password",
			"dev_reason": err.Error(),
		})
		return
	}

	err = user.UpdatePassword(newPassword)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to update password",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// ResetUserPassword godoc
//
//	@Summary		Reset User Password
//	@Description	Reset user password with otp
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			password	body		models.RESET_PASSWORD		true	"User Password & Otp"
//	@Success		200			{object}	models.HTTP_USER_RESPONSE	"Password Reset successfully"
//	@Failure		400			{object}	models.Error
//	@Failure		404			{object}	models.Error
//	@Failure		409			{object}	models.Error
//	@Failure		417			{object}	models.Error
//	@Failure		500			{object}	models.Error
//	@Router			/user/reset-password [patch]
func ResetUserPassword(context *gin.Context) {
	var resetPasswordData models.RESET_PASSWORD

	err := context.ShouldBindJSON(&resetPasswordData)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	user, err := models.GetUserByEmail(resetPasswordData.Email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "The email is not registered",
			"dev_reason": err.Error(),
		})
		return
	}

	otp, err := db.RDB.Get(db.Ctx, resetPasswordData.Email).Result()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message":    "OTP expired",
			"dev_reason": err.Error(),
		})
		return
	}

	if otp != resetPasswordData.OTP {
		context.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message":    "Invalid OTP",
			"dev_reason": "OTP does not match",
		})
		return
	}

	hashPassword, err := utils.HashText(resetPasswordData.Password)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error occured while securing password",
			"dev_reason": err.Error(),
		})
		return
	}

	err = user.UpdatePassword(hashPassword)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"message":    "Error occured while saving password",
			"dev_reason": err.Error(),
		})

		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})

}

// ResetUserPin godoc
//
//	@Summary		Reset User Pin
//	@Description	Reset user pin with otp
//	@Tags			User
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			pin	body		models.RESET_PIN			true	"User Pin & Otp"
//	@Success		200	{object}	models.HTTP_USER_RESPONSE	"Pin Reset successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		409	{object}	models.Error
//	@Failure		417	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/user/reset-pin [patch]
func ResetUserPin(context *gin.Context) {
	var requestPinData models.RESET_PIN
	email := context.GetString("email")
	err := context.ShouldBindJSON(&requestPinData)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}
	otp, err := db.RDB.Get(db.Ctx, email).Result()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message":    "OTP expired",
			"dev_reason": err.Error(),
		})
		return
	}

	if otp != requestPinData.OTP {
		context.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message":    "Invalid OTP",
			"dev_reason": "OTP does not match",
		})
		return
	}

	hashPin, err := utils.HashText(requestPinData.Pin)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error occured while securing password",
			"dev_reason": err.Error(),
		})
		return
	}

	err = user.UpdatePin(hashPin)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error occured while saving pin",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Pin reset successfully",
	})

}
