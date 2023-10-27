package bot

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	errorFormat = "*Error:* _%s_"
)

// ReportError sends back a formatted error message to the channel where we received the event from
func (msgCtx *SlackMessageContext) ReportError(err error) {

	apiClient := msgCtx.api
	event := msgCtx.innerEvent

	opts := []slack.MsgOption{
		slack.MsgOptionText(fmt.Sprintf(errorFormat, err.Error()), false),
	}

	opts = append(opts, slack.MsgOptionTS(event.TimeStamp))

	_, _, err = apiClient.PostMessage(msgCtx.channelID, opts...)
	if err != nil {
		fmt.Printf("failed posting message: %v\n", err)
	}
}

// Reply send a message to the current channel
func (msgCtx *SlackMessageContext) Reply(message string, options ...slack.MsgOption) error {
	ev := msgCtx.innerEvent
	if ev == nil {
		return fmt.Errorf("unable to get message event details")
	}
	return msgCtx.Post(msgCtx.channelID, message, options...)
}

// Post send a message to a channel
func (msgCtx *SlackMessageContext) Post(channel string, message string, options ...slack.MsgOption) error {

	apiClient := msgCtx.api
	event := msgCtx.innerEvent
	if event == nil {
		return fmt.Errorf("unable to get message event details")
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
	}

	_, _, err := apiClient.PostMessage(
		channel,
		opts...,
	)
	return err
}
