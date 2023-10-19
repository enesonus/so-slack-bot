package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

func TestDatabase(database *Queries, count int) {

	start := time.Now()
	for i := 0; i < count; i++ {

		params := CreateBotParams{
			CreatedAt:      time.Now(),
			LastActivityAt: time.Now(),
			BotToken:       "firstBotToken",
		}

		_, err := database.CreateBot(context.Background(), params)
		if err != nil {
			log.Fatal("Couldn't create bot: ", err)
		}
	}
	fmt.Printf("Time to create %v bots: %v\n", count, time.Since(start))

	start = time.Now()
	bots, err := database.GetBots(context.Background())
	fmt.Printf("Time to get bots: %v\n", time.Since(start))

	if err != nil {
		log.Fatal("Couldn't get bots: ", err)
	}
	// fmt.Println("Bot just created: ", botFromDB)
	fmt.Printf("Total Bot count: %v\n", len(bots))
}
