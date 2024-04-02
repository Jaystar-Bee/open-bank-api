package models

import (
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

type TRANSACTION struct {
	ID              int64     `json:"id"`
	Sender          int64     `json:"sender" binding:"required"`
	Sender_Wallet   int64     `json:"sender_wallet" binding:"required"`
	Receiver        int64     `json:"receiver" binding:"required"`
	Receiver_Wallet int64     `json:"receiver_wallet" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	Status          string    `json:"status" binding:"required"`
	Remarks         string    `json:"remarks"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

func (transaction *TRANSACTION) Save() (*TRANSACTION, error) {
	query := `INSERT INTO transactions (sender, sender_wallet, receiver, receiver_wallet, amount, status, type, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return transaction, err
	}
	defer statement.Close()
	data, err := statement.Exec(transaction.Sender, transaction.Sender_Wallet, transaction.Receiver, transaction.Receiver_Wallet, transaction.Amount, transaction.Status, transaction.Type, utils.NowTime(), nil, nil)
	if err != nil {
		return transaction, err
	}
	transaction.ID, err = data.LastInsertId()
	return transaction, err
}

func (transaction *TRANSACTION) Update() (*TRANSACTION, error) {
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
