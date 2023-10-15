package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/shomali11/slacker"
)

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

	return database
}

func StartSlackBot(slackBotToken string) {

	fmt.Println(slackBotToken)
	slackBot := slacker.NewClient(slackBotToken, os.Getenv("SLACK_APP_TOKEN"))

	slackBot.Command("set_so_channel", setSOChannelDef)
	slackBot.Command("remove_so_channel", removeSOChannelDef)
	slackBot.Command("getinfo", getUserInfoDef)
	slackBot.Command("add_tag {tag}", setSOChannelDef)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go PrintCommandEvents(slackBot.CommandEvents())

	var err error
	go func() {
		err = slackBot.Listen(ctx)
	}()
	if err != nil {
		log.Printf("Error listening to Slack Bot: %v, Token: %v\n", err, slackBotToken)
	}
	botParams := db.CreateBotParams{
		BotToken:       slackBotToken,
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}
	databaseObject := setDatabase()
	databaseObject.CreateBot(context.Background(), botParams)
	fmt.Println("Bot is not listening...")
}
