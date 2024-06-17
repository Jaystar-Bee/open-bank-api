package models

type Upload struct {
	File string `json:"file" binding:"required"`
	Id   int    `json:"id" binding:"required"`
}
