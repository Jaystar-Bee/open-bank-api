package models

import (
	"errors"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

type USER struct {
	ID                   int64     `json:"id"`
	FirstName            string    `json:"first_name" binding:"required"`
	LastName             string    `json:"last_name" binding:"required"`
	Email                string    `json:"email" binding:"required"`
	Password             string    `json:"password" binding:"required"`
	Phone                string    `json:"phone"`
	TransactionPin       string    `json:"transaction_pin" binding:"required"`
	Tag                  string    `json:"tag" binding:"required"`
	IsVerified           bool      `json:"is_verified"`
	Avatar               string    `json:"avatar"`
	AccountIsDeactivated bool      `json:"account_is_deactivated"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeletedAt            time.Time `json:"deleted_at"`
}

type USER_LOGIN struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp" binding:"required"`
}

type OTP_REQUEST struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name"`
}

type USER_RESPONSE struct {
	ID                   int64      `json:"id"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	Email                string     `json:"email"`
	Phone                string     `json:"phone"`
	Tag                  string     `json:"tag"`
	IsVerified           bool       `json:"is_verified"`
	Avatar               *string    `json:"avatar"`
	AccountIsDeactivated bool       `json:"account_is_deactivated"`
	CreatedAt            *time.Time `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at"`
}

func (user *USER) Save() error {
	query := `INSERT INTO users (first_name, last_name, email, password, phone, transaction_pin, tag, is_verified, account_is_deactivated, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

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

	err = db.MainDB.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.TransactionPin, user.Tag, false, false, utils.NowTime()).Scan(&user.ID)
	if err != nil {
		return err
	}
	err = user.CreateWallet()
	if err != nil {
		return errors.New("error creating wallet")

	}
	return nil
}

func (login *USER_LOGIN) Login() error {
	query := `SELECT password FROM users WHERE email = $1`

	var password string
	err := db.MainDB.QueryRow(query, login.Email).Scan(&password)
	if err != nil {
		return err
	}

	err = utils.CompareHash(password, login.Password)
	return err
}

func GetUserByEmail(email string) (*USER_RESPONSE, error) {
	query := `SELECT id, first_name, last_name, email, phone, tag, is_verified, avatar, account_is_deactivated, created_at, updated_at, deleted_at FROM users WHERE email = $1`
	data := db.MainDB.QueryRow(query, email)

	user := &USER_RESPONSE{}
	err := data.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Tag, &user.IsVerified, &user.Avatar, &user.AccountIsDeactivated, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	return user, err

}
func GetUserByPhone(phone string) (*USER_RESPONSE, error) {
	query := `SELECT id, first_name, last_name, email, phone, tag, is_verified, avatar, account_is_deactivated, created_at, updated_at, deleted_at FROM users WHERE phone = $1`
	data := db.MainDB.QueryRow(query, phone)

	user := &USER_RESPONSE{}
	err := data.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Tag, &user.IsVerified, &user.Avatar, &user.AccountIsDeactivated, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	return user, err

}
func GetUserByTag(tag string) (*USER_RESPONSE, error) {
	query := `SELECT id, first_name, last_name, email, phone, tag, is_verified, avatar, account_is_deactivated, created_at, updated_at, deleted_at FROM users WHERE tag = $1`
	data := db.MainDB.QueryRow(query, tag)

	user := &USER_RESPONSE{}
	err := data.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Tag, &user.IsVerified, &user.Avatar, &user.AccountIsDeactivated, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	return user, err
}

func GetUserByID(id int64) (*USER_RESPONSE, error) {
	query := `SELECT id, first_name, last_name, email, phone, tag, is_verified, avatar, account_is_deactivated, created_at, updated_at, deleted_at FROM users WHERE id = $1`
	data := db.MainDB.QueryRow(query, id)

	user := &USER_RESPONSE{}
	err := data.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Tag, &user.IsVerified, &user.Avatar, &user.AccountIsDeactivated, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	return user, err
}

func (user *USER_RESPONSE) ConfirmPin(pin string) error {
	query := `SELECT transaction_pin FROM users WHERE id = $1`
	var dbPin string
	err := db.MainDB.QueryRow(query, user.ID).Scan(&dbPin)
	if err != nil {
		return err
	}
	err = utils.CompareHash(dbPin, pin)
	if err != nil {
		return err
	}
	return nil
}

func (user *USER_RESPONSE) ConfirmPassword(password string) error {
	query := `SELECT password FROM users WHERE id = $1`
	var dbPassword string
	err := db.MainDB.QueryRow(query, user.ID).Scan(&dbPassword)
	if err != nil {
		return err
	}
	err = utils.CompareHash(dbPassword, password)
	if err != nil {
		return err
	}
	return nil
}

func (user *USER_RESPONSE) UpdateUser() error {
	query := `UPDATE users SET first_name = $1, last_name = $2, email = $3, phone = $4, tag = $5, is_verified = $6, avatar = $7, account_is_deactivated = $8, updated_at = $9 WHERE id = $10`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(user.FirstName, user.LastName, user.Email, user.Phone, user.Tag, user.IsVerified, user.Avatar, user.AccountIsDeactivated, utils.NowTime(), user.ID)
	if err != nil {
		return err
	}
	return nil
}
func (user *USER_RESPONSE) UpdatePassword(hashPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(hashPassword, utils.NowTime(), user.ID)
	return err
}

func (user *USER_RESPONSE) UpdatePin(hashPin string) error {
	query := `UPDATE users SET transaction_pin = $1, updated_at = $2 WHERE id = $3`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(hashPin, utils.NowTime(), user.ID)
	return err
}
