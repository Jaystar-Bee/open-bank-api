package models

type USER_REQUEST struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Phone          string `json:"phone"`
	TransactionPin string `json:"transaction_pin" binding:"required"`
	Tag            string `json:"tag" binding:"required"`
	Avatar         string `json:"avatar"`
}

type DEPOSIT_BODY struct {
	Amount float64 `json:"amount" binding:"required"`
}

type HTTP_REQUEST struct {
	ID        int64   `json:"id"`
	Requester int64   `json:"requester"`
	Giver     int64   `json:"giver" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Status    string  `json:"status"`
	Remarks   string  `json:"remarks"`
}
