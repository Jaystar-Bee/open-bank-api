package models

import (
	"database/sql"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

type REQUEST struct {
	ID        int64          `json:"id"`
	Requester int64          `json:"requester"`
	Giver     int64          `json:"giver" binding:"required"`
	Amount    float64        `json:"amount" binding:"required"`
	Status    string         `json:"status"`
	Remarks   string         `json:"remarks"`
	CreatedAt sql.NullString `json:"created_at"`
	UpdatedAt sql.NullString `json:"updated_at"`
	DeletedAt sql.NullString `json:"deleted_at"`
}

func (request *REQUEST) Save() (*REQUEST, error) {
	query := `INSERT INTO requests (requester, giver, amount, status, remarks, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := db.MainDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	data, err := stmt.Exec(request.Requester, request.Giver, request.Amount, request.Status, request.Remarks, utils.NowTime())
	if err != nil {
		return nil, err
	}
	request.ID, err = data.LastInsertId()
	return request, err
}

func (request *REQUEST) Update() (*REQUEST, error) {
	query := `UPDATE requests SET status = $1, updated_at = $2 WHERE id = $3`
	stmt, err := db.MainDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(request.Status, utils.NowTime(), request.ID)
	if err != nil {
		return nil, err
	}
	return request, err
}
func (request *REQUEST) Delete() error {
	query := `DELETE FROM requests WHERE id = $1`
	stmt, err := db.MainDB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(request.ID)
	if err != nil {
		return err
	}
	return err
}

func GetRequestByID(requestId int64) (*REQUEST, error) {
	query := `SELECT * FROM requests WHERE id = $1`
	row := db.MainDB.QueryRow(query, requestId)

	request := &REQUEST{}
	err := row.Scan(&request.ID, &request.Requester, &request.Giver, &request.Amount, &request.Status, &request.Remarks, &request.CreatedAt, &request.UpdatedAt, &request.DeletedAt)

	return request, err
}
