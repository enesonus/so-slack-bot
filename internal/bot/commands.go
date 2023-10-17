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

func StartSlackBot(slackBotToken string) error {

	fmt.Println("Slack Bot Token: ", slackBotToken)
	slackBot := slacker.NewClient(slackBotToken, os.Getenv("SLACK_APP_TOKEN"))

	slackBot.Command("set_so_channel", setSOChannelDef)
	slackBot.Command("remove_so_channel", removeSOChannelDef)
	slackBot.Command("getinfo", getUserInfoDef)
	slackBot.Command("add_tag {tag}", addTagDef)

	ctx, cancel := context.WithCancel(context.Background())

	go PrintCommandEvents(slackBot.CommandEvents())

	err := error(nil)
	go func() {
		err = slackBot.Listen(ctx)
		defer cancel()
	}()
	if err != nil {
		log.Printf("Error listening to Slack Bot: %v, Token: %v\n", err, slackBotToken)
		return fmt.Errorf("error listening to slack bot: %v, token: %v", err, slackBotToken)
	}

	teamInfo, err := slackBot.APIClient().GetTeamInfo()
	if err != nil {
		log.Printf("Error getting team info: %v\n", err)
		return fmt.Errorf("error getting team info: %v", err)
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
		log.Printf("Error creating workspace: %v\n", err)
		return fmt.Errorf("error creating workspace: %v", err)
	}
	botParams := db.CreateBotParams{
		BotToken:       slackBotToken,
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		WorkspaceID:    teamInfo.ID,
	}
	_, err = databaseObject.CreateBot(context.Background(), botParams)
	if err != nil {
		log.Printf("Error creating bot: %v\n", err)
		return fmt.Errorf("error creating bot: %v", err)
	}

	return nil
}
