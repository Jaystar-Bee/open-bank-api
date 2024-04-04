package models

type USER_REQUEST struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Phone          string `json:"phone"`
	TransactionPin string `json:"transaction_pin" binding:"required"`
	Tag            string `json:"tag" binding:"required"`
}
