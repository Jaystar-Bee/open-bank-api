package db

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

var MainDB *sql.DB
var Ctx = context.Background()
var RDB *redis.Client

func InitDatabase() {
	initRedis()
	Database, err := sql.Open("sqlite3", "api.db")
	if err != nil {
		panic(err)
	}

	MainDB = Database
	MainDB.SetMaxOpenConns(14)
	MainDB.SetMaxIdleConns(7)
	CreateTables()
}

func initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	RDB = rdb

}
