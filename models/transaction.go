package models

import (
	"database/sql"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

const (
	Transaction_send    = "DEBIT"
	Transaction_receive = "CREDIT"
	Transaction_return  = "REVERTED"
)
const (
	Transaction_pending   = "PENDING"
	Transaction_failed    = "FAILED"
	Transaction_rejected  = "REJECTED"
	Transaction_completed = "COMPLETED"
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
}

type TRANSACTION_RESPONSE[T int64 | *USER_RESPONSE] struct {
	ID              int64          `json:"id"`
	Sender          T              `json:"sender" binding:"required"`
	Sender_Wallet   int64          `json:"sender_wallet" binding:"required"`
	Receiver        T              `json:"receiver" binding:"required"`
	Receiver_Wallet int64          `json:"receiver_wallet" binding:"required"`
	Amount          float64        `json:"amount" binding:"required"`
	Status          string         `json:"status" binding:"required"`
	Remarks         sql.NullString `json:"remarks"`
	Type            string         `json:"type"`
	CreatedAt       sql.NullString `json:"created_at"`
	UpdatedAt       sql.NullString `json:"updated_at"`
	DeletedAt       sql.NullString `json:"deleted_at"`
}

func (transaction *TRANSACTION[int64]) Save() (*TRANSACTION[int64], error) {
	query := `INSERT INTO transactions (sender, sender_wallet, receiver, receiver_wallet, amount, status, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return transaction, err
	}
	defer statement.Close()
	data, err := statement.Exec(transaction.Sender, transaction.Sender_Wallet, transaction.Receiver, transaction.Receiver_Wallet, transaction.Amount, transaction.Status, utils.NowTime(), nil, nil)
	if err != nil {
		return transaction, err
	}
	transaction.ID, err = data.LastInsertId()
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
		err := rows.Scan(&transaction.ID, &sender, &transaction.Sender_Wallet, &receiver, &transaction.Receiver_Wallet, &transaction.Amount, &transaction.Status, &transaction.Remarks, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt)
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
	err := db.MainDB.QueryRow(query, id).Scan(&transaction.ID, &sender, &transaction.Sender_Wallet, &receiver, &transaction.Receiver_Wallet, &transaction.Amount, &transaction.Status, &transaction.Remarks, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt)
	if err != nil {
		return nil, err
	}
	transaction.Sender, _ = GetUserByID(sender)
	transaction.Receiver, _ = GetUserByID(receiver)
	return &transaction, nil
}
