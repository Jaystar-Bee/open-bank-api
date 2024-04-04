package models

import (
	"database/sql"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

type WALLET_REQUEST struct {
	ID        int64     `json:"id"`
	Balance   float64   `json:"balance"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type WALLET struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	Balance   float64        `json:"balance"`
	CreatedAt sql.NullString `json:"created_at"`
	UpdatedAt sql.NullString `json:"updated_at"`
	DeletedAt sql.NullString `json:"deleted_at"`
}

type ADD_TO_BALANCE_BODY struct {
	ID             int64   `json:"id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Remarks        string  `json:"remarks"`
	TransactionPin string  `json:"transaction_pin" binding:"required"`
}

func (user *USER) CreateWallet() error {
	query := "INSERT INTO wallets (user_id, balance, created_at) VALUES (?, ?, ?)"
	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID, 120000, utils.NowTime())
	return err
}

func (user *USER_RESPONSE) GetWallet() (*WALLET, error) {
	query := `
	SELECT * FROM wallets WHERE user_id = ?
	`
	wallet := &WALLET{}
	err := db.MainDB.QueryRow(query, user.ID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	if err != nil {
		return nil, err
	}
	return wallet, nil

}

func AddToBalance(amount float64, userId int64) error {
	query := `
	UPDATE wallets SET balance = balance + ?, updated_at = ? WHERE user_id = ?
	`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(amount, utils.NowTime(), userId)
	return err
}

func (wallet *WALLET) RemoveFromBalance(amount float64, transaction *TRANSACTION[int64]) (*TRANSACTION[int64], error) {
	query := `UPDATE wallets SET balance = balance - ?, updated_at = ? WHERE id = ?`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return transaction, err
	}

	defer statement.Close()
	_, err = statement.Exec(amount, utils.NowTime(), wallet.ID)
	if err != nil {
		return transaction, err
	}
	transaction, err = transaction.Save()
	return transaction, err
}
