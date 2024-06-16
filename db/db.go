package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

var MainDB *sql.DB
var Ctx = context.Background()
var RDB *redis.Client

func InitDatabase() {
	initRedis()
	if os.Getenv("ENV") == "production" {
		initPostgres()
	} else {
		initDev()
	}
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
func initPostgres() {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	MainDB = db
}

func initDev() {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	MainDB = db
}
