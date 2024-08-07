package models

import (
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/utils"
)

const (
	Request_Giver     = "GIVER"
	Request_Requester = "REQUESTER"
)

type REQUEST struct {
	ID        int64      `json:"id"`
	Requester int64      `json:"requester"`
	Giver     int64      `json:"giver" binding:"required"`
	Amount    float64    `json:"amount" binding:"required"`
	Status    string     `json:"status"`
	Remarks   string     `json:"remarks"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type REQUEST_RESPONSE struct {
	ID        int64         `json:"id"`
	Requester USER_RESPONSE `json:"requester"`
	Giver     USER_RESPONSE `json:"giver"`
	Amount    float64       `json:"amount"`
	Status    string        `json:"status"`
	Remarks   string        `json:"remarks"`
	CreatedAt *time.Time    `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at"`
}

func (request *REQUEST) Save() (*REQUEST, error) {
	query := `INSERT INTO requests (requester, giver, amount, status, remarks, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := db.MainDB.QueryRow(query, request.Requester, request.Giver, request.Amount, request.Status, request.Remarks, utils.NowTime()).Scan(&request.ID)
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

func GetUserResquests(userId int64) ([]*REQUEST, error) {
	query := `SELECT * FROM requests WHERE requester = $1`
	rows, err := db.MainDB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*REQUEST
	for rows.Next() {
		request := &REQUEST{}
		err := rows.Scan(&request.ID, &request.Requester, &request.Giver, &request.Amount, &request.Status, &request.Remarks, &request.CreatedAt, &request.UpdatedAt, &request.DeletedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func GetRequestsToPay(userId int64) ([]*REQUEST, error) {
	query := `SELECT * FROM requests WHERE giver = $1`
	rows, err := db.MainDB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*REQUEST
	for rows.Next() {
		request := &REQUEST{}
		err := rows.Scan(&request.ID, &request.Requester, &request.Giver, &request.Amount, &request.Status, &request.Remarks, &request.CreatedAt, &request.UpdatedAt, &request.DeletedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func GetAllRequests(userId int64) ([]*REQUEST, error) {
	query := `SELECT * FROM requests WHERE requester = $1 OR giver = $1`
	rows, err := db.MainDB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*REQUEST
	for rows.Next() {
		request := &REQUEST{}
		err := rows.Scan(&request.ID, &request.Requester, &request.Giver, &request.Amount, &request.Status, &request.Remarks, &request.CreatedAt, &request.UpdatedAt, &request.DeletedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}
