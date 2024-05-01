package models

type HTTP_USER_RESPONSE struct {
	Data    USER_RESPONSE `json:"data"`
	Message string        `json:"message"`
}

type HTTP_TOKEN_RESPONSE struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type HTTP_LOGIN_RESPONSE struct {
	Data    USER   `json:"data"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type HTTP_MESSAGE_ONLY_RESPONSE struct {
	Message string `json:"message"`
}

type HTTP_TRANSACTION_BY_ID_RESPONSE struct {
	Data    TRANSACTION[*USER_RESPONSE] `json:"data"`
	Message string                      `json:"message"`
}

type HTTP_TRANSACTION_DATA_RESPONSE struct {
	Transactions  TRANSACTION[*USER_RESPONSE] `json:"transactions"`
	Total_counts  float64                     `json:"total_counts"`
	Page_number   float64                     `json:"page_number"`
	Per_page      float64                     `json:"per_page"`
	Total_pages   float64                     `json:"total_pages"`
	Current_page  float64                     `json:"current_page"`
	Next_page     float64                     `json:"next_page"`
	Previous_page float64                     `json:"previous_page"`
	First_page    float64                     `json:"first_page"`
	Last_page     float64                     `json:"last_page"`
	Has_next      bool                        `json:"has_next"`
	Has_previous  bool                        `json:"has_previous"`
}
type HTTP_TRANSACTION_LIST_RESPONSE struct {
	Data    HTTP_TRANSACTION_DATA_RESPONSE `json:"data"`
	Message string                         `json:"message"`
}

type HTTP_WALLET_RESPONSE struct {
	Data    WALLET_REQUEST `json:"data"`
	Message string         `json:"message"`
}

type HTTP_REQUEST_RESPONSE struct {
	Data    HTTP_REQUEST `json:"data"`
	Message string       `json:"message"`
}
