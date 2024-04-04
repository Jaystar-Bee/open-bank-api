package models

type Error struct {
	Message   string `json:"message"`
	DevReason string `json:"dev_reason"`
}
