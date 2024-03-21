package models

import (
	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

type WALLET struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt string  `json:"deleted_at"`
}

func (user USER) CreateWallet() error {
	query := "INSERT INTO wallets (user_id, balance, created_at) VALUES (?, ?, ?)"
	statement, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID, 0, utils.NowTime())
	return err
}
