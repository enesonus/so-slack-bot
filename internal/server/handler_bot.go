package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/enesonus/so-slack-bot/internal/bot"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func EventsAPIHandler(w http.ResponseWriter, r *http.Request) {
	body, err := GetBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validateRequest(w, r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Error at validateRequest: %v\n", err)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error at ParseEvent: %v", err)
		return
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
		switch innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			msgCtx, err := bot.SlackMessageHandler(w, &eventsAPIEvent)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			go msgCtx.Listen("show_tags", bot.ShowTags)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func validateRequest(w http.ResponseWriter, r *http.Request) error {

	body, err := GetBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error at GetBody: %v", err)
	}

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("SLACK_SIGNING_SECRET must be set")
	}

	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
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

func GetBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	// Read the body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}
