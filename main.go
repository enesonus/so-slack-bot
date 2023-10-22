package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/enesonus/so-slack-bot/internal/bot"
	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/enesonus/so-slack-bot/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func startAllBots() {
	fmt.Println("Starting all bots")
	botCount := 0

	databaseObject, err := db.GetDatabase()
	if err != nil {
		log.Fatal("Couldn't get database: ", err)
	}

	bots, err := databaseObject.GetBots(context.Background())
	if err != nil {
		log.Fatal("Couldn't get bots: ", err)
	}

	for _, newbot := range bots {
		go func(token string) { _, err = bot.StartSlackBot(token) }(newbot.BotToken)
		botCount++
		if err != nil {
			fmt.Printf("Error starting bot: %v\n", err)
			botCount--
		}
	}

	fmt.Printf("%v bots in DB, %v bots started\n", len(bots), botCount)
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
	go bot.QuestionCheckerAndSender()
	srv.ListenAndServe()

}
