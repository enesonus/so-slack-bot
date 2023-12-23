package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/enesonus/so-slack-bot/internal/bot"
	"github.com/enesonus/so-slack-bot/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

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
	router.Post("/events-api-handler", server.EventsAPIHandler)

	go bot.QuestionCheckerAndSender()
	srv.ListenAndServe()

}
