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

type USER_EDIT struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Phone     string  `json:"phone"`
	Tag       string  `json:"tag" binding:"required"`
	Avatar    *string `json:"avatar"`
}

type DEPOSIT_BODY struct {
	Amount float64 `json:"amount" binding:"required"`
}

type CHANGE_PIN struct {
	OldPin string `json:"old_pin" binding:"required"`
	NewPin string `json:"new_pin" binding:"required"`
}

type RESET_PASSWORD struct {
	OTP      string `json:"otp" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RESET_PIN struct {
	OTP string `json:"otp" binding:"required"`
	Pin string `json:"pin" binding:"required"`
}
type CHANGE_PASSWORD struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type HTTP_REQUEST struct {
	ID        int64   `json:"id"`
	Requester int64   `json:"requester"`
	Giver     int64   `json:"giver" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Status    string  `json:"status"`
	Remarks   string  `json:"remarks"`
}

type HTTP_FILE_REQUEST struct {
	File string `json:"file" binding:"required"`
	Id   string `json:"id" binding:"required"`
}
