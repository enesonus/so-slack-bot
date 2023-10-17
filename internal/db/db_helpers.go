package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDatabase() (*Queries, error) {

	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Print("DATABASE_URL must be set")
		return nil, fmt.Errorf("DATABASE_URL must be set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print("Can't connect to DB: ", err)
	}

	database := New(conn)

	return database, err
}
