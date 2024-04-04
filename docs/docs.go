// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "John Ayilara (Jaystar)",
            "url": "https://github.com/Jaystar-Bee",
            "email": "jbayilara@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/transactions": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "You can get user transaction list and the list are paginated, which is 10 transactions per page by default.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get user transaction list",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_TRANSACTION_LIST_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Check queries",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unable to fetch transactions",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/transactions/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "You can get transaction by ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get transaction by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Transaction ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_TRANSACTION_BY_ID_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Check queries",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Transaction not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Unable to fetch transaction",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "Log user in to the application.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Log user in",
                "parameters": [
                    {
                        "description": "Log User In",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.USER_LOGIN"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User logged in successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_LOGIN_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/renew": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Renew user token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Renew token",
                "responses": {
                    "200": {
                        "description": "Token renewd successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_TOKEN_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/signup": {
            "post": {
                "description": "Onboard user to the application.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Create a user",
                "parameters": [
                    {
                        "description": "Create User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.USER_REQUEST"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_MESSAGE_ONLY_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Check body",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/{email}": {
            "get": {
                "description": "Get user by email.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user by email.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Email",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_USER_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/{phone}": {
            "get": {
                "description": "Get user by phone.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user by phone.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Phone",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_USER_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/{tag}": {
            "get": {
                "description": "Get user by tag.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user by tag.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Tag",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_USER_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/user/{user_id}": {
            "get": {
                "description": "Get user by Id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user by Id.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Id",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_USER_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/wallet": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get Wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Get Wallet",
                "responses": {
                    "200": {
                        "description": "wallet fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_WALLET_RESPONSE"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/wallet/send": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Send money to another user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Send money",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ADD_TO_BALANCE_BODY"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "wallet updated successfully",
                        "schema": {
                            "$ref": "#/definitions/models.HTTP_TRANSACTION_BY_ID_RESPONSE"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ADD_TO_BALANCE_BODY": {
            "type": "object",
            "required": [
                "amount",
                "id",
                "transaction_pin"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "remarks": {
                    "type": "string"
                },
                "transaction_pin": {
                    "type": "string"
                }
            }
        },
        "models.Error": {
            "type": "object",
            "properties": {
                "dev_reason": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_LOGIN_RESPONSE": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.USER"
                },
                "message": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_MESSAGE_ONLY_RESPONSE": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_TOKEN_RESPONSE": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_TRANSACTION_BY_ID_RESPONSE": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.TRANSACTION-models_USER"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_TRANSACTION_DATA_RESPONSE": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "number"
                },
                "first_page": {
                    "type": "number"
                },
                "has_next": {
                    "type": "boolean"
                },
                "has_previous": {
                    "type": "boolean"
                },
                "last_page": {
                    "type": "number"
                },
                "next_page": {
                    "type": "number"
                },
                "page_number": {
                    "type": "number"
                },
                "per_page": {
                    "type": "number"
                },
                "previous_page": {
                    "type": "number"
                },
                "total_counts": {
                    "type": "number"
                },
                "total_pages": {
                    "type": "number"
                },
                "transactions": {
                    "$ref": "#/definitions/models.TRANSACTION-models_USER"
                }
            }
        },
        "models.HTTP_TRANSACTION_LIST_RESPONSE": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.HTTP_TRANSACTION_DATA_RESPONSE"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_USER_RESPONSE": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.USER"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.HTTP_WALLET_RESPONSE": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.WALLET_REQUEST"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.TRANSACTION-models_USER": {
            "type": "object",
            "required": [
                "amount",
                "receiver",
                "receiver_wallet",
                "sender",
                "sender_wallet",
                "status"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "receiver": {
                    "$ref": "#/definitions/models.USER"
                },
                "receiver_wallet": {
                    "type": "integer"
                },
                "remarks": {
                    "type": "string"
                },
                "sender": {
                    "$ref": "#/definitions/models.USER"
                },
                "sender_wallet": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.USER": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "password",
                "tag",
                "transaction_pin"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_verified": {
                    "type": "boolean"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "tag": {
                    "type": "string"
                },
                "transaction_pin": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.USER_LOGIN": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.USER_REQUEST": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "password",
                "tag",
                "transaction_pin"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "tag": {
                    "type": "string"
                },
                "transaction_pin": {
                    "type": "string"
                }
            }
        },
        "models.WALLET_REQUEST": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "OPEN BANK API",
	Description:      "OPEN BANK API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
