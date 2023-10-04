package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func main() {
	botInit()

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Command("set_so_channel", setSOChannelDef)
	bot.Command("remove_so_channel", removeSOChannelDef)
	bot.Command("getinfo", getUserInfoDef)
	bot.Command("add_tag {tag}", setSOChannelDef)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go printCommandEvents(bot.CommandEvents())

	// go botStackOverflow(bot, "open-telemetry")

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
