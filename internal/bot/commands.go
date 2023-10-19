package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/shomali11/slacker"
)

func StartSlackBot(slackBotToken string) (*slacker.Slacker, error) {

	slackBot := slacker.NewClient(slackBotToken, os.Getenv("SLACK_APP_TOKEN"))
	slackBot.Command("set_so_channel", setSOChannelDef)
	slackBot.Command("remove_so_channel", removeSOChannelDef)
	slackBot.Command("getinfo", getUserInfoDef)
	slackBot.Command("show_tags", showTagsDef)
	slackBot.Command("add_tag {tag}", addTagDef)

	go PrintCommandEvents(slackBot.CommandEvents())

	// Add bot and Workspace to DB

	teamInfo, err := slackBot.APIClient().GetTeamInfo()
	if err != nil {
		log.Printf("Error getting team info: %v\n", err)
		return nil, fmt.Errorf("error getting team info: %v", err)
	}

	workspaceParams := db.GetOrCreateWorkspaceParams{
		ID:              teamInfo.ID,
		WorkspaceName:   teamInfo.Name,
		WorkspaceDomain: teamInfo.Domain,
		CreatedAt:       time.Now(),
	}

	databaseObject, err := db.GetDatabase()
	if err != nil {
		fmt.Printf("Error connecting database: %v\n", err)
	}

	_, err = databaseObject.GetOrCreateWorkspace(context.Background(), workspaceParams)
	if err != nil {
		log.Printf("Error adding workspace to DB: %v\n", err)
		return nil, fmt.Errorf("error adding workspace to DB: %v", err)
	}
	botParams := db.CreateBotParams{
		BotToken:       slackBotToken,
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		WorkspaceID:    teamInfo.ID,
	}
	_, err = databaseObject.CreateBot(context.Background(), botParams)
	if err != nil {
		log.Printf("Error adding bot to DB: %v\n", err)
		return slackBot, fmt.Errorf("error adding bot to DB: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err = slackBot.Listen(ctx)
		defer cancel()
	}()

	if err != nil {
		log.Printf("Error listening to Slack Bot: %v, Token: %v\n", err, slackBotToken)
		return nil, fmt.Errorf("error listening to slack bot: %v, token: %v", err, slackBotToken)
	}

	return slackBot, nil
}
