package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func EventsAPIHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validateRequest(w, body, r.Header); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Error at validateRequest: %v\n", err)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error at ParseEvent: %v", err)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var res *slackevents.ChallengeResponse
		err = json.Unmarshal([]byte(body), &res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(res.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if ev.BotID != "" {
				w.WriteHeader(http.StatusOK)
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("Error at GetBotByWorkspaceID: %v\n", err)
				return
			}
			api, err := getSlackAPIClient(eventsAPIEvent.TeamID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("Error at getSlackAPIClient: %v\n", err)
				return
			}
			fmt.Printf("Message: %v\n", ev.Text)
			w.WriteHeader(http.StatusOK)
			api.PostMessage(ev.Channel, slack.MsgOptionText(
				fmt.Sprintf("Message: %v\n", ev.Text), false))
		}
	}
}

func validateRequest(w http.ResponseWriter, body []byte, header http.Header) error {

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("SLACK_SIGNING_SECRET must be set")
	}

	sv, err := slack.NewSecretsVerifier(header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("error at NewSecretsVerifier: %v", err)
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error at sv.Write: %v", err)
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("error at sv.Ensure: %v", err)
	}
	return nil
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
