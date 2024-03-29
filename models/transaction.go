package models

import (
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

const (
	Transaction_send    = "SEND"
	Transaction_receive = "RECEIVED"
	Transaction_return  = "REVERTED"
)
const (
	Transaction_pending   = "PENDING"
	Transaction_failed    = "FAILED"
	Transaction_completed = "COMPLETED"
)

type TRANSACTION struct {
	ID              int64     `json:"id"`
	Sender          int64     `json:"sender" binding:"required"`
	Sender_Wallet   int64     `json:"sender_wallet" binding:"required"`
	Receiver        int64     `json:"receiver" binding:"required"`
	Receiver_Wallet int64     `json:"receiver_wallet" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	Status          bool      `json:"status" binding:"required"`
	Type            bool      `json:"type" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

func (transaction *TRANSACTION) Save() (*TRANSACTION, error) {
	query := `INSERT INTO transactions (sender, sender_wallet, receiver, receiver_wallet, amount, status, type, created_at, updated_at, deleted_at) VALUES (%1, %2, %3, %4, %5, %6, %7, %8, %9, %10)`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	data, err := statement.Exec(transaction.Sender, transaction.Sender_Wallet, transaction.Receiver, transaction.Receiver_Wallet, transaction.Amount, transaction.Status, transaction.Type, utils.NowTime(), nil, nil)
	if err != nil {
		return nil, err
	}

	transaction.ID, err = data.LastInsertId()
	return transaction, err
}
