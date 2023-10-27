package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func PrintCommandEvents(slackChannel <-chan *slacker.CommandEvent) {
	for event := range slackChannel {
		log.Printf("Command Event Received")
		log.Printf("Command: %v", event.Command)
		log.Printf("Parameters: %v", event.Parameters)
		log.Printf("Event: %v\n\n", event.Event)
	}
}

func CreateSlackBot(slackBotToken string) (*slacker.Slacker, error) {

	slackBot := slacker.NewClient(slackBotToken, os.Getenv("SLACK_APP_TOKEN"))
	slackBot.Command("set_so_channel", setSOChannelDef)
	slackBot.Command("remove_so_channel", removeSOChannelDef)
	// slackBot.Command("getinfo", getUserInfoDef)
	slackBot.Command("add_tag {tag}", addTagDef)
	// slackBot.Command("show_tags", showTagsDef)

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
	defer cancel()

	err = slackBot.Listen(ctx)

	if err != nil {
		log.Printf("Error listening to Slack Bot: %v, Token: %v\n", err, slackBotToken)
		return nil, fmt.Errorf("error listening to slack bot: %v, token: %v", err, slackBotToken)
	}

	return slackBot, nil
}

func StartSlackBot(slackBotToken string) (*slacker.Slacker, error) {
	fmt.Printf("Starting slack bot with token: %v\n", slackBotToken)
	slackBot := slacker.NewClient(slackBotToken, os.Getenv("SLACK_APP_TOKEN"))
	slackBot.Command("set_so_channel", setSOChannelDef)
	slackBot.Command("remove_so_channel", removeSOChannelDef)
	// slackBot.Command("getinfo", getUserInfoDef)
	slackBot.Command("add_tag {tag}", addTagDef)
	// slackBot.Command("show_tags", showTagsDef)

	go PrintCommandEvents(slackBot.CommandEvents())

	err := error(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = slackBot.Listen(ctx)
	if err != nil {
		log.Printf("Error listening to Slack Bot: %v, Token: %v\n", err, slackBotToken)
		return nil, fmt.Errorf("error listening to slack bot: %v, token: %v", err, slackBotToken)
	}

	return slackBot, nil
}

func SlackMessageHandler(w http.ResponseWriter, eventsAPIEvent *slackevents.EventsAPIEvent) (*SlackMessageContext, error) {
	start := time.Now()
	innerEvent := eventsAPIEvent.InnerEvent
	ev := innerEvent.Data.(*slackevents.MessageEvent)
	if ev.BotID != "" {
		w.WriteHeader(http.StatusOK)
		return &SlackMessageContext{}, nil
	}
	msgCtx, err := NewClient(eventsAPIEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return &SlackMessageContext{}, fmt.Errorf("error at NewClient: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("SlackMessageHandler took %v\n", time.Since(start))
	return msgCtx, nil
}

func getSlackAPIClient(workspaceID string) (*slack.Client, error) {

	dbObj, err := db.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("error at GetDatabase: %v", err)
	}
	botObj, err := dbObj.GetBotByWorkspaceID(context.Background(), workspaceID)
	if err != nil {
		return nil, fmt.Errorf("error at GetBotByWorkspaceID: %v", err)
	}
	api := slack.New(botObj.BotToken)
	return api, nil
}

func (messageCtx *SlackMessageContext) ChannelName() (string, error) {
	channel, err := messageCtx.api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID:         messageCtx.channelID,
		IncludeLocale:     false,
		IncludeNumMembers: false})
	if err != nil {
		return "", fmt.Errorf("error getting channel: %v", err)
	}
	return channel.Name, nil
}

type CommandFunc func(msgCtx *SlackMessageContext, suffix string)
type SlackMessageContext struct {
	api            *slack.Client
	eventsAPIEvent *slackevents.EventsAPIEvent
	innerEvent     *slackevents.MessageEvent
	message        string
	prefix         string
	channelID      string
}

func NewClient(eventsAPIEvent *slackevents.EventsAPIEvent) (*SlackMessageContext, error) {

	innerEvent := eventsAPIEvent.InnerEvent
	if innerEvent.Type != "message" {
		return &SlackMessageContext{}, fmt.Errorf("error at NewClient: innerEvent.Type is not MessageEvent")
	}
	ev := innerEvent.Data.(*slackevents.MessageEvent)
	message := ev.Text
	words := strings.Fields(message)
	prefix := ""
	if len(words) > 0 {
		prefix = words[0]
	}
	api, err := getSlackAPIClient(eventsAPIEvent.TeamID)
	if err != nil {
		return &SlackMessageContext{}, fmt.Errorf("error at NewClient/getSlackAPIClient: %v", err)
	}
	return &SlackMessageContext{
		api:            api,
		eventsAPIEvent: eventsAPIEvent,
		innerEvent:     ev,
		message:        message,
		prefix:         prefix,
		channelID:      ev.Channel,
	}, nil
}

func (messageCtx *SlackMessageContext) Listen(commandPrefix string, commandFunc CommandFunc) {
	suffix := ""
	if len(messageCtx.message) > len(commandPrefix) {
		suffix = messageCtx.message[len(commandPrefix):]
	}
	if messageCtx.prefix == commandPrefix {
		commandFunc(messageCtx, suffix)
	}
}
