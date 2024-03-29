package handlers

import (
	"net/http"

	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/gin-gonic/gin"
)

func GetWallet(context *gin.Context) {

	user_id := context.GetInt64("user")

	user, err := models.GetUserByID(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not Found",
			"dev_reason": err.Error(),
		})
		return
	}

	wallet, err := user.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Failed to get user wallet",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Wallet fetched successfully",
		"data":    wallet,
	})

}

func SendToUser(context *gin.Context) {

	var body models.ADD_TO_BALANCE_BODY

	err := context.ShouldBindJSON(&body)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unabale to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	user_id := context.GetInt64("user")

	sender, err := models.GetUserByID(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Sender not Found",
			"dev_reason": err.Error(),
		})
		return
	}
	receiver, err := models.GetUserByID(body.ID)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Receiver not Found",
			"dev_reason": err.Error(),
		})
		return
	}

	sender_wallet, err := sender.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Unable to get sender's wallet",
			"dev_reason": err.Error(),
		})
	}
	if sender_wallet.Balance == 0 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You have no money, Please top up your wallet",
		})
		return
	}
	if sender_wallet.Balance < float64(body.Amount) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You do not have enough for this transaction, Please try lower amount",
		})
		return
	}

	err = sender_wallet.RemoveFromBalance(float64(body.Amount))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fulfill transaction",
			"dev_reason": err.Error(),
		})
		return
	}
	err = models.AddToBalance(float64(body.Amount), body.ID)
	if err != nil {
		models.AddToBalance(float64(body.Amount), sender.ID)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fulfill transaction",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Transaction successfull",
	})

}
