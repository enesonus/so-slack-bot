package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/enesonus/so-slack-bot/internal/bot"
	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/enesonus/so-slack-bot/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

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

func startAllBots() {

	databaseObject, err := db.GetDatabase()
	if err != nil {
		log.Fatal("Couldn't get database: ", err)
	}
	bots, err := databaseObject.GetBots(context.Background())
	if err != nil {
		log.Fatal("Couldn't get bots: ", err)
	}
	for _, newbot := range bots {
		bot.StartSlackBot(newbot.BotToken)
	}
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	fmt.Println("Port: ", port)

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
	router.Get("/access_token/", server.GetAccessTokenAndStartBot)

	startAllBots()
	srv.ListenAndServe()

}
