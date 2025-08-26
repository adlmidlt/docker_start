package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var DB *sql.DB

func Init() {
	dsn := fmt.Sprintf("host=postgresql port=5432 user=user1 password=pass123word dbname=demo sslmode=disable")

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database")
}

func GetDB() *sql.DB {
	return DB
}
