package models

import (
	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

type WALLET struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt string  `json:"deleted_at"`
}

type ADD_TO_BALANCE_BODY struct {
	ID     int64 `json:"id" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
}

func (user *USER) CreateWallet() error {
	query := "INSERT INTO wallets (user_id, balance, created_at) VALUES (?, ?, ?)"
	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID, 0, utils.NowTime())
	return err
}

func (user *USER_RESPONSE) GetWallet() (*WALLET, error) {
	query := `
	SELECT * FROM wallets WHERE user_id = ? AND deleted_at IS NULL
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
	UPDATE wallets SET balance = balance + ? WHERE user_id = ?
	`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(amount, userId)
	return err
}

func (wallet *WALLET) RemoveFromBalance(amount float64) error {
	query := `UPDATE wallets SET balance = balance - ? WHERE id = ?`

	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(amount, wallet.ID)
	return err
}
