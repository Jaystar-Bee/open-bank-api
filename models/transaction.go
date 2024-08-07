package models

import (
	"database/sql"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

const (
	Transaction_send      = "DEBIT"
	Transaction_receive   = "CREDIT"
	Transaction_return    = "REVERTED"
	Transaction_cancelled = "CANCELLED"
)
const (
	Transaction_pending   = "PENDING"
	Transaction_failed    = "FAILED"
	Transaction_rejected  = "REJECTED"
	Transaction_completed = "COMPLETED"
)

const (
	Transaction_Channel_Request = "REQUEST"
	Transaction_Channel_Wallet  = "WALLET"
)

type UserTypes interface {
}

type TRANSACTION[T int64 | *USER_RESPONSE | *USER] struct {
	ID              int64     `json:"id"`
	Sender          T         `json:"sender" binding:"required"`
	Sender_Wallet   int64     `json:"sender_wallet" binding:"required"`
	Receiver        T         `json:"receiver" binding:"required"`
	Receiver_Wallet int64     `json:"receiver_wallet" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	Status          string    `json:"status" binding:"required"`
	Remarks         string    `json:"remarks"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
	Channel         string    `json:"channel"`
	// Channel         string    `json:"channel" binding:"required"` // TODO: Change to thi one. Once done, this is required
}

type TRANSACTION_RESPONSE[T int64 | *USER_RESPONSE] struct {
	ID              int64      `json:"id"`
	Sender          T          `json:"sender" binding:"required"`
	Sender_Wallet   int64      `json:"sender_wallet" binding:"required"`
	Receiver        T          `json:"receiver" binding:"required"`
	Receiver_Wallet int64      `json:"receiver_wallet" binding:"required"`
	Amount          float64    `json:"amount" binding:"required"`
	Status          string     `json:"status" binding:"required"`
	Remarks         string     `json:"remarks"`
	Type            string     `json:"type"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	Channel         string     `json:"channel"`
	// Channel         string    `json:"channel" binding:"required"` // TODO: Change to thi one. Once done, this is required
}

func (transaction *TRANSACTION[int64]) Save() (*TRANSACTION[int64], error) {
	query := `INSERT INTO transactions (sender, sender_wallet, receiver, receiver_wallet, amount, status, channel, remarks, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	err := db.MainDB.QueryRow(query, transaction.Sender, transaction.Sender_Wallet, transaction.Receiver, transaction.Receiver_Wallet, transaction.Amount, transaction.Status, transaction.Channel, transaction.Remarks, utils.NowTime(), nil, nil).Scan(&transaction.ID)
	return transaction, err
}

func (transaction *TRANSACTION[int64]) Update() (*TRANSACTION[int64], error) {
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return transaction, err
	}
	defer statement.Close()

	_, err = statement.Exec(transaction.Status, utils.NowTime(), transaction.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (user *USER_RESPONSE) GetTransactions(per_page, page_number float64) ([]TRANSACTION_RESPONSE[*USER_RESPONSE], float64, error) {

	query := `SELECT * FROM transactions WHERE sender = $1 OR receiver = $1 LIMIT $2 OFFSET $3`

	queryCount := `SELECT COUNT(*) FROM transactions WHERE sender = $1 OR receiver = $1`

	var count float64
	err := db.MainDB.QueryRow(queryCount, user.ID).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.MainDB.Query(query, user.ID, per_page, per_page*(page_number-1))
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var transactions []TRANSACTION_RESPONSE[*USER_RESPONSE]
	for rows.Next() {
		var transaction TRANSACTION_RESPONSE[*USER_RESPONSE]
		var sender int64
		var receiver int64
		err := rows.Scan(&transaction.ID, &sender, &transaction.Sender_Wallet, &receiver, &transaction.Receiver_Wallet, &transaction.Amount, &transaction.Status, &transaction.Remarks, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt, &transaction.Channel)
		if err != nil {
			return nil, 0, err
		}
		transaction.Sender, _ = GetUserByID(sender)
		transaction.Receiver, _ = GetUserByID(receiver)
		transactions = append(transactions, transaction)
	}
	return transactions, count, nil
}

func GetTransactionByID(id int64) (*TRANSACTION_RESPONSE[*USER_RESPONSE], error) {
	query := `SELECT * FROM transactions WHERE id = $1`

	var transaction TRANSACTION_RESPONSE[*USER_RESPONSE]
	var sender int64
	var receiver int64
	err := db.MainDB.QueryRow(query, id).Scan(&transaction.ID, &sender, &transaction.Sender_Wallet, &receiver, &transaction.Receiver_Wallet, &transaction.Amount, &transaction.Status, &transaction.Remarks, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt, &transaction.Channel)
	if err != nil {
		return nil, err
	}
	transaction.Sender, _ = GetUserByID(sender)
	transaction.Receiver, _ = GetUserByID(receiver)
	return &transaction, nil
}

func GetMoneyIn(userId int64) (float64, error) {
	query := `SELECT SUM(amount) FROM transactions WHERE receiver = $1 AND status = $2`

	var amount sql.NullFloat64
	err := db.MainDB.QueryRow(query, userId, Transaction_completed).Scan(&amount)
	if err != nil {
		return 0, err
	}
	if amount.Valid {
		return amount.Float64, nil
	}
	return 0, nil
}

func GetMoneyOut(userId int64) (float64, error) {
	query := `SELECT SUM(amount) FROM transactions WHERE sender = $1 AND status = $2`

	var amount sql.NullFloat64
	err := db.MainDB.QueryRow(query, userId, Transaction_completed).Scan(&amount)
	if err != nil {
		return 0, err
	}
	if amount.Valid {
		return amount.Float64, nil
	}
	return 0, nil
}
