package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/enesonus/go-slack-bot/internal/db"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to DB: ", err)
	}

	db := db.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	fmt.Printf("%v\n", apiCfg)

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
