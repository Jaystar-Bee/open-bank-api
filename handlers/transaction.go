package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/gin-gonic/gin"
)

// GetTags 		godoc
//
//	@Summary		Get user transaction list
//	@Description	You can get user transaction list and the list are paginated, which is 10 transactions per page by default.
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page_number	query		int									false	"Page number"
//	@Param			per_page	query		int									false	"Transactions per page"
//	@Success		200			{object}	models.HTTP_TRANSACTION_LIST_RESPONSE	"Ok"
//	@Failure		400			{object}	models.Error							"Check queries"
//	@Failure		404			{object}	models.Error							"User not found"
//	@Failure		500			{object}	models.Error							"Unable to fetch transactions"
//	@Router			/transactions [get]
func GetTransactions(context *gin.Context) {
	user_id := context.GetInt64("user")

	user, err := models.GetUserByID(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}

	page_number_string, ok := context.GetQuery("page_number")
	if !ok {
		page_number_string = "1"
	}
	page_number, err := strconv.ParseFloat(page_number_string, 64)

	if err != nil || page_number < 1 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}
	per_page_string, ok := context.GetQuery("per_page")
	if !ok {
		per_page_string = "10"
	}
	per_page, err := strconv.ParseFloat(per_page_string, 64)
	if err != nil || per_page < 1 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	transactions, total_counts, err := user.GetTransactions(per_page, page_number)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to fetch transactions",
			"dev_reason": err.Error(),
		})
		return
	}

	for key, transaction := range transactions {
		if transaction.Receiver.ID == user_id {
			transactions[key].Type = models.Transaction_receive
		} else {
			transactions[key].Type = models.Transaction_send
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data": gin.H{
			"transactions":  transactions,
			"total_counts":  total_counts,
			"page_number":   page_number,
			"per_page":      per_page,
			"total_pages":   math.Ceil(total_counts / per_page),
			"current_page":  page_number,
			"next_page":     page_number + 1,
			"previous_page": page_number - 1,
			"first_page":    1,
			"last_page":     math.Ceil(total_counts / per_page),
			"has_next":      page_number < math.Ceil(total_counts/per_page),
			"has_previous":  page_number > 1,
		},
	})

}

// GetTransactionByID 	godoc
//
//	@Tags			Transactions
//	@Summary		Get transaction by ID
//	@Description	You can get transaction by ID.
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int										true	"Transaction ID"
//	@Success		200	{object}	models.HTTP_TRANSACTION_BY_ID_RESPONSE	"Ok"
//	@Failure		400	{object}	models.Error							"Check queries"
//	@Failure		404	{object}	models.Error							"User not found"
//	@Failure		404	{object}	models.Error							"Transaction not found"
//	@Failure		500	{object}	models.Error							"Unable to fetch transaction"
//	@Router			/transactions/{id} [get]
func GetTransactionByID(context *gin.Context) {
	user_id := context.GetInt64("user")

	_, err := models.GetUserByID(user_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "User not found",
			"dev_reason": err.Error(),
		})
		return
	}

	transaction_id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":    "Unable to process request",
			"dev_reason": err.Error(),
		})
		return
	}

	transaction, err := models.GetTransactionByID(transaction_id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message":    "Unable to fetch transaction",
			"dev_reason": err.Error(),
		})
		return
	}
	if transaction.Receiver.ID != user_id && transaction.Sender.ID != user_id {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "You're not authorized to this transaction",
		})
		return
	}

	if transaction.Receiver.ID == user_id {
		transaction.Type = models.Transaction_receive
	} else {
		transaction.Type = models.Transaction_send
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Transaction fetched successfully",
		"data":    transaction,
	})

}
