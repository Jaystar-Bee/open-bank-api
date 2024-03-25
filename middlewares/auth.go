package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/jwt"
	"github.com/gin-gonic/gin"
)

func CheckAuthentication(context *gin.Context) {
	token := context.GetHeader("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}

	token = strings.Split(token, "Bearer ")[1]

	claimsData, err := jwt.ValidateJWT(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
	}

	expireTime := int64(claimsData["ExpiredAt"].(float64))

	if expireTime < time.Now().Unix() {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "session timeout"})
	}

	context.Set("user", claimsData["user"])
	context.Next()
}
