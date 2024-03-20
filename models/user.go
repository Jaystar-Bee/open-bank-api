package models

import (
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

type USER struct {
	ID             int64     `json:"id"`
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	Password       string    `json:"password" binding:"required"`
	Phone          string    `json:"phone" binding:"required"`
	TransactionPin string    `json:"transaction_pin"`
	AccountNumber  string    `json:"account_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

type USER_LOGIN struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type USER_RESPONSE struct {
	ID            int64     `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	AccountNumber string    `json:"account_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (user USER) Save() error {
	query := `INSERT INTO users (first_name, last_name, email, password, phone, transaction_pin, account_number, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}

	hashPassword, err := utils.HashText(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	hashPin, err := utils.HashText(user.TransactionPin)
	if err != nil {
		return err
	}
	user.TransactionPin = hashPin

	defer statement.Close()
	_, err = statement.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.TransactionPin, user.AccountNumber, utils.NowTime(), nil, nil)
	return err
}

func GetUserByEmail(email string) (*USER_RESPONSE, error) {
	query := `SELECT * FROM users WHERE email = $1`
	data := db.MainDB.QueryRow(query, email)

	user := &USER_RESPONSE{}
	err := data.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.AccountNumber, &user.CreatedAt, &user.UpdatedAt)
	return user, err

}
