package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/jwt"
	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/gin-gonic/gin"
)

func CheckAuthentication(context *gin.Context) {
	token := context.GetHeader("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	token = strings.Split(token, "Bearer ")[1]

	claimsData, err := jwt.ValidateJWT(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}
	expireTime := int64(claimsData["ExpiredAt"].(float64))
	if expireTime < time.Now().Unix() {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "session timeout"})
		return
	}
	user := int64(claimsData["userId"].(float64))
	email := claimsData["email"].(string)

	context.Set("user", user)
	context.Set("email", email)
	context.Next()
}

func CheckAccountActivation(context *gin.Context) {
	userId := context.GetInt64("user")

	user, err := models.GetUserByID(userId)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not Found", "dev_reason": err.Error()})
		return
	}
	if user.AccountIsDeactivated {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Account is deactivated"})
		return
	}

	context.Next()
}
