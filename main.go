package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/enesonus/so-slack-bot/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func setDatabase() *db.Queries {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to DB: ", err)
	}

	database := db.New(conn)
	apiCfg := apiConfig{
		DB: database,
	}
	fmt.Printf("apiCfg: %v\n", apiCfg)

	return database
}

func testDatabase(database *db.Queries, count int) {

	start := time.Now()
	for i := 0; i < count; i++ {

		params := db.CreateBotParams{
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

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	fmt.Println("Port: ", port)

	databaseObject := setDatabase()
	server.FetchTagsAndSaveToDB(databaseObject, 1, 654)

	// testDatabase(databaseObject, 5)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"https:*", "http:*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			MaxAge:         300,
		}))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	router.Get("/", server.CheckReadiness)
	router.Get("/healthz", server.CheckReadiness)
	router.Get("/access_token/", server.GetAccessToken)

	srv.ListenAndServe()

}
