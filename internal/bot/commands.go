package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// func PrintCommandEvents(slackChannel <-chan *slacker.CommandEvent) {
// 	for event := range slackChannel {
// 		log.Printf("Command Event Received")
// 		log.Printf("Command: %v", event.Command)
// 		log.Printf("Parameters: %v", event.Parameters)
// 		log.Printf("Event: %v\n\n", event.Event)
// 	}
// }

func CreateSlackBot(slackBotToken string) (*slack.Client, error) {

	// Add bot and Workspace to DB

	apiClient := slack.New(slackBotToken)

	teamInfo, err := apiClient.GetTeamInfo()
	if err != nil {
		log.Printf("Error getting team info: %v\n", err)
		return &slack.Client{}, fmt.Errorf("error getting team info: %v", err)
	}

	workspaceParams := db.GetOrCreateWorkspaceParams{
		ID:              teamInfo.ID,
		WorkspaceName:   teamInfo.Name,
		WorkspaceDomain: teamInfo.Domain,
		CreatedAt:       time.Now(),
	}

	if err != nil {
		fmt.Printf("Error connecting database: %v\n", err)
	}

	_, err = dbObj.GetOrCreateWorkspace(context.Background(), workspaceParams)
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
	_, err = dbObj.CreateBot(context.Background(), botParams)
	if err != nil {
		log.Printf("Error adding bot to DB: %v\n", err)
		return apiClient, fmt.Errorf("error adding bot to DB: %v", err)
	}

	return apiClient, nil
}

var dbObj, err = db.GetDatabase()

// ticker := time.NewTicker(time.Duration(timePeriod) * time.Minute)

func NewMessageContext(w http.ResponseWriter, eventsAPIEvent *slackevents.EventsAPIEvent) (*SlackMessageContext, error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("can not connect to DB: %v", err)
		dbObj, err = db.GetDatabase()
	}
	start := time.Now()
	msgCtx, err := NewClient(eventsAPIEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return &SlackMessageContext{}, fmt.Errorf("error at NewClient: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("NewMessageContext took %v\n", time.Since(start))
	return msgCtx, nil
}

func IsBot(ev *slackevents.MessageEvent) bool {
	return ev.BotID != ""
}

func getSlackAPIClient(workspaceID string) (*slack.Client, error) {

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

func (msgCtx *SlackMessageContext) ChannelName() (string, error) {
	channelName := ""
	channel, err := dbObj.GetChannelByID(context.Background(), msgCtx.ChannelID)
	channelName = channel.ChannelName
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			slackChannel, err := msgCtx.Api.GetConversationInfo(&slack.GetConversationInfoInput{
				ChannelID:         msgCtx.ChannelID,
				IncludeLocale:     false,
				IncludeNumMembers: false})
			if err != nil {
				return "", fmt.Errorf("error getting channel: %v", err)
			}
			channelName = slackChannel.Name
		} else {
			return "", fmt.Errorf("error getting channel: %v", err)
		}
	}
	return channelName, nil
}

type CommandFunc func(msgCtx *SlackMessageContext, suffix string)
type SlackMessageContext struct {
	Api            *slack.Client
	EventsAPIEvent *slackevents.EventsAPIEvent
	InnerEvent     *slackevents.MessageEvent
	Message        string
	Prefix         string
	ChannelID      string
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
		Api:            api,
		EventsAPIEvent: eventsAPIEvent,
		InnerEvent:     ev,
		Message:        message,
		Prefix:         prefix,
		ChannelID:      ev.Channel,
	}, nil
}

func (messageCtx *SlackMessageContext) Command(commandPrefix string, commandFunc CommandFunc) {
	suffix := ""
	words := strings.Fields(messageCtx.Message)
	if len(words) > 1 {
		suffix = strings.Join(words[1:], "")
	}
	if messageCtx.Prefix == commandPrefix {
		commandFunc(messageCtx, suffix)
	}
}
