package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/gin-gonic/gin"
)

// GetWallet godoc
//
//	@Summary		Get Wallet
//	@Description	Get Wallet
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	models.HTTP_WALLET_RESPONSE	"wallet fetched successfully"
//	@Failure		401	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Router			/wallet [get]
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

// TopUpWallet godoc
//
//	@Summary		Send money
//	@Description	Send money to another user
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		models.ADD_TO_BALANCE_BODY				true	"body"
//	@Success		200		{object}	models.HTTP_TRANSACTION_BY_ID_RESPONSE	"wallet updated successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		401		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Router			/wallet/send [post]
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

	if user_id == body.ID {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You cannot send money to yourself",
		})
		return
	}

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

	receiver_wallet, err := receiver.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Unable to get receiver's wallet",
			"dev_reason": err.Error(),
		})
	}

	err = sender.ConfirmPin(body.TransactionPin)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Invalid transaction pin",
			"dev_reason": err.Error(),
		})
		return
	}

	senderTransaction := &models.TRANSACTION[int64]{
		Sender:          sender.ID,
		Sender_Wallet:   sender_wallet.ID,
		Receiver:        receiver.ID,
		Receiver_Wallet: receiver_wallet.ID,
		Amount:          body.Amount,
		Status:          models.Transaction_pending,
		Channel:         models.Transaction_Channel_Wallet,
		Remarks:         body.Remarks,
		CreatedAt:       time.Now(),
	}
	// Type:            models.Transaction_send,

	senderTransaction, err = sender_wallet.RemoveFromBalance(float64(body.Amount), senderTransaction)
	if err != nil {
		senderTransaction.Status = models.Transaction_failed
		senderTransaction.Update()
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fulfill transaction",
			"data":       senderTransaction,
			"dev_reason": err.Error(),
		})
		return
	}
	err = models.AddToBalance(float64(body.Amount), body.ID)
	if err != nil {
		senderTransaction.Status = models.Transaction_failed
		senderTransaction.Update()
		models.AddToBalance(float64(body.Amount), sender.ID)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fulfill transaction",
			"data":       senderTransaction,
			"dev_reason": err.Error(),
		})
		return
	}
	senderTransaction.Status = models.Transaction_completed
	senderTransaction.Update()
	context.JSON(http.StatusCreated, gin.H{
		"message": "Transaction successfull",
		"data":    senderTransaction,
	})

}

// Deposit godoc
//
//	@Summary		Deposit money
//	@Description	Deposit money to your wallet
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		models.DEPOSIT_BODY						true	"body"
//	@Success		200		{object}	models.HTTP_TRANSACTION_BY_ID_RESPONSE	"wallet updated successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		401		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Router			/wallet/deposit [post]
func Deposit(context *gin.Context) {
	var body struct {
		Amount float64 `json:"amount"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	if body.Amount <= 0 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Amount must be greater than 0",
		})
		return
	}
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
			"message":    "Unable to get user wallet",
			"dev_reason": err.Error(),
		})
		return
	}
	transaction := &models.TRANSACTION[int64]{
		Sender:          user.ID,
		Sender_Wallet:   wallet.ID,
		Amount:          body.Amount,
		Receiver:        user.ID,
		Receiver_Wallet: wallet.ID,
		Status:          models.Transaction_pending,
		Channel:         models.Transaction_Channel_Wallet,
		Remarks:         "Deposit",
		CreatedAt:       time.Now(),
	}
	transaction, err = transaction.Save()
	if err != nil {
		transaction.Status = models.Transaction_failed
		transaction.Update()
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fulfill transaction",
			"data":       transaction,
			"dev_reason": err.Error(),
		})
		return
	}
	err = models.AddToBalance(body.Amount, user.ID)
	if err != nil {
		transaction.Status = models.Transaction_failed
		transaction.Update()
		models.AddToBalance(body.Amount, user.ID)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fulfill transaction",
			"data":       transaction,
			"dev_reason": err.Error(),
		})
		return
	}

	transaction.Status = models.Transaction_completed
	transaction.Update()
	context.JSON(http.StatusCreated, gin.H{
		"message": "Transaction successfull",
		"data":    transaction,
	})

}

// RequestMoney	godoc
//
//	@Summary		Request money
//	@Description	Request money from another user
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		models.HTTP_REQUEST					true	"body"
//	@Success		200		{object}	models.HTTP_MESSAGE_ONLY_RESPONSE	"Request successfully"
//	@Failure		400		{object}	models.Error
//	@Failure		401		{object}	models.Error
//	@Failure		404		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Router			/wallet/requests [post]
func RequestMoney(context *gin.Context) {
	var body models.REQUEST
	err := context.ShouldBindJSON(&body)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	if body.Amount <= 0 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Amount must be greater than 0",
		})
		return
	}
	user_id := context.GetInt64("user")
	user, err := models.GetUserByID(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not Found",
			"dev_reason": err.Error(),
		})
		return
	}
	_, err = models.GetUserByID(body.Giver)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Can't find the other user",
			"dev_reason": err.Error(),
		})
		return
	}
	if user.ID == body.Giver {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You cannot request money from yourself",
		})
		return
	}
	body.Requester = user.ID
	body.Status = models.Transaction_pending

	_, err = body.Save()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to make a request",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Request successfull",
	})

}

// DeleteRequest godoc
//
//	@Summary		Delete a request
//	@Description	Delete a request
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Request ID"
//	@Success		200	{object}	models.HTTP_MESSAGE_ONLY_RESPONSE
//	@Failure		400	{object}	models.Error
//	@Failure		401	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/wallet/requests/{id} [delete]
func DeleteRequest(context *gin.Context) {

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	request, err := models.GetRequestByID(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Request not found",
			"dev_reason": err.Error(),
		})
		return
	}

	if request.Status != models.Transaction_pending {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request can't be deleted",
		})
		return
	}

	if request.Requester != context.GetInt64("user") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You are not authorized to delete this request",
		})
		return
	}

	err = request.Delete()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to delete request",
			"dev_reason": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Request deleted successfully",
	})

}

// GetUserRequest godoc
//
//	@Summary		Get request
//	@Description	Get the money request you made and also have. You can differentiate the get by adding type query to be GIVER OR REQUESTER
//	@Accept			json
//	@Produce		json
//	@Tags			Wallet
//	@Security		ApiKeyAuth
//	@Param			type				query		string	false	"GIVER OR REQUESTER"
//	@Failure		400					{object}	models.Error
//	@Failure		401					{object}	models.Error
//	@Failure		404					{object}	models.Error
//	@Failure		500					{object}	models.Error
//	@Success		200					{object}	models.HTTP_REQUEST_RESPONSE	"Request fetched successfully"
//	@Router			/wallet/requests	[get]
func GetUserRequests(context *gin.Context) {
	user_id := context.GetInt64("user")

	requestType := context.Query("type") // GIVER OR REQUESTER

	_, err := models.GetUserByID(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not Found",
			"dev_reason": err.Error(),
		})
		return
	}
	var requests []*models.REQUEST
	if requestType == models.Request_Giver {
		requests, err = models.GetRequestsToPay(user_id)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message":    "Unable to get requests",
				"dev_reason": err.Error(),
			})
			return
		}
	} else {
		requests, err = models.GetUserResquests(user_id)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message":    "Unable to get requests",
				"dev_reason": err.Error(),
			})
			return
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Requests fetched successfully",
		"data":    requests,
		"count":   len(requests),
	})

}

// RejectUserRequest godoc
//
//	@Summary		Reject a user money request
//	@Description	Reject a user money request
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int								true	"Request ID"
//	@Success		200	{object}	models.HTTP_REQUEST_RESPONSE	"Request rejected successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		401	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/wallet/requests/{id}/reject [post]
func RejectRequest(context *gin.Context) {

	// PARSE PARAMS
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	request, err := models.GetRequestByID(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Request not found",
			"dev_reason": err.Error(),
		})
		return
	}

	// CHECKS
	if request.Status != models.Transaction_pending {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request can't be rejected",
		})
		return
	}

	if request.Giver != context.GetInt64("user") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You are not authorized to reject this request",
		})
		return
	}

	// GET GIVER AND REQUESTER WALLETS
	giver, err := models.GetUserByID(request.Giver)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error rejecting request",
			"dev_reason": err.Error(),
		})
		return
	}
	giver_wallet, err := giver.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error rejecting request",
			"dev_reason": err.Error(),
		})
		return
	}
	requester, err := models.GetUserByID(request.Requester)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error rejecting request",
			"dev_reason": err.Error(),
		})
		return
	}
	requester_wallet, err := requester.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error rejecting request",
			"dev_reason": err.Error(),
		})
		return
	}

	// INITIATE TRANSACTION
	transaction := &models.TRANSACTION[int64]{
		Sender:          request.Giver,
		Sender_Wallet:   giver_wallet.ID,
		Amount:          request.Amount,
		Receiver:        request.Requester,
		Receiver_Wallet: requester_wallet.ID,
		Status:          models.Transaction_pending,
		Channel:         models.Transaction_Channel_Request,
		Remarks:         request.Remarks,
		CreatedAt:       time.Now(),
	}
	transaction, err = transaction.Save()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error rejecting request",
			"dev_reason": err.Error(),
		})
		return
	}

	// SET REQUEST STATUS
	request.Status = models.Transaction_rejected
	request, err = request.Update()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error rejecting request",
			"dev_reason": err.Error(),
		})
		return
	}
	transaction.Status = models.Transaction_rejected
	_, _ = transaction.Update()

	context.JSON(http.StatusOK, gin.H{
		"message": "Request rejected successfully",
		"data":    request,
	})

}

// AcceptUserRequest godoc
//
//	@Summary		Accept a user money request
//	@Description	Accept a user money request
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int								true	"Request ID"
//	@Success		200	{object}	models.HTTP_REQUEST_RESPONSE	"Request accepted successfully"
//	@Failure		400	{object}	models.Error
//	@Failure		401	{object}	models.Error
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/wallet/requests/{id}/confirm [post]
func ConfirmRequest(context *gin.Context) {

	// PARSE PARAMS
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	request, err := models.GetRequestByID(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Request not found",
			"dev_reason": err.Error(),
		})
		return
	}

	// CHECKS
	if request.Status != models.Transaction_pending {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request can't be confirmed",
		})
		return
	}

	if request.Giver != context.GetInt64("user") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You are not authorized to confirm this request",
		})
		return
	}

	// GET GIVER AND REQUESTER WALLETS
	giver, err := models.GetUserByID(request.Giver)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}
	giver_wallet, err := giver.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}

	if request.Amount > giver_wallet.Balance {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Insufficient balance",
		})
		return
	}

	requester, err := models.GetUserByID(request.Requester)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}
	requester_wallet, err := requester.GetWallet()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}

	// INITIATE TRANSACTION
	transaction := &models.TRANSACTION[int64]{
		Sender:          request.Giver,
		Sender_Wallet:   giver_wallet.ID,
		Amount:          request.Amount,
		Receiver:        request.Requester,
		Receiver_Wallet: requester_wallet.ID,
		Status:          models.Transaction_pending,
		Channel:         models.Transaction_Channel_Request,
		Remarks:         request.Remarks,
		CreatedAt:       time.Now(),
	}

	// REMOVE FROM GIVER BALANCE AND SAVE TRANSACTION AS PENDING
	transaction, err = giver_wallet.RemoveFromBalance(request.Amount, transaction)
	if err != nil {
		transaction.Status = models.Transaction_failed
		transaction.Update()
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}

	// ADD TO REQUESTER BALANCE
	err = models.AddToBalance(request.Amount, requester_wallet.ID)
	if err != nil {
		transaction.Status = models.Transaction_failed
		transaction.Update()
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}

	// UPDATE REQUEST TO SUCCESS
	request.Status = models.Transaction_completed
	request, err = request.Update()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}
	transaction.Status = models.Transaction_completed
	_, err = transaction.Update()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Error confirming request",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Request confirmed successfully",
		"data":    request,
	})

}
