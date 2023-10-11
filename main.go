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
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	fmt.Println("Port: ", port)

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

	start := time.Now()
	for i := 0; i < 1000; i++ {

		params := db.CreateBotParams{
			ID:             uuid.UUID.String(uuid.New()),
			CreatedAt:      time.Now(),
			LastActivityAt: time.Now(),
			BotToken:       "firstBotToken",
		}

		_, err := database.CreateBot(context.Background(), params)
		if err != nil {
			log.Fatal("Couldn't create bot: ", err)
		}
	}
	fmt.Printf("Time to create 1 bots: %v\n", time.Since(start))

	start = time.Now()
	bots, err := database.GetBots(context.Background())
	fmt.Printf("Time to get bots: %v\n", time.Since(start))

	if err != nil {
		log.Fatal("Couldn't get bots: ", err)
	}
	// fmt.Println("Bot just created: ", botFromDB)
	fmt.Printf("Bot count: %v\n", len(bots))

	router := chi.NewRouter()

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

	go srv.ListenAndServe()

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Command("set_so_channel", setSOChannelDef)
	bot.Command("remove_so_channel", removeSOChannelDef)
	bot.Command("getinfo", getUserInfoDef)
	bot.Command("add_tag {tag}", setSOChannelDef)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go printCommandEvents(bot.CommandEvents())

	// go botStackOverflow(bot, "open-telemetry")

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
