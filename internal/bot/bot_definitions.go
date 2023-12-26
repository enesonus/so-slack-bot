package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/slack-go/slack"
)

// ##############################################################################################################################
// bot.Command("upload <sentence>", uploadDef)
// bot.Command("ping", pingPong)
// bot.Command("echo {word} {word2}", echoWordDef)
// bot.Command("say <sentence>", saySentenceDef)
// bot.Command("repeat <word> {number}", repeatNtimesDef)
// bot.Command("message", messageReplyDefinition)

func RemoveSOChannelDef(msgCtx *SlackMessageContext, suffix string) {
	apiClient := msgCtx.Api

	if msgCtx.ChannelID != "" {
		_, err := dbObj.DeleteChannel(context.Background(), msgCtx.ChannelID)
		if err != nil {
			if strings.Contains(err.Error(), "rows in result set"){
				apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
					"Channel is already not a Stack Overflow Notification channel", false))
					return
			}

			log.Printf("Error deleting channel: %v\n", err)
			return
		}
		channelName, _ := msgCtx.ChannelName()
		fmt.Printf("Stack Overflow channel is removed\n")
		apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
			fmt.Sprintf("`%s` is removed from notification channels", channelName), false))
	}
}

func SetSOChannelDef(msgCtx *SlackMessageContext, suffix string) {
	apiClient := msgCtx.Api

	if msgCtx.ChannelID != "" {
		team, err := apiClient.GetTeamInfo()
		if err != nil {
			fmt.Printf("Error getting team info: %v\n", err)
		}
		if err != nil {
			fmt.Printf("Error connecting database: %v\n", err)
		}

		bot, err := dbObj.GetBotByWorkspaceID(context.Background(), team.ID)
		if err != nil {
			fmt.Printf("Error getting bot: %v\n", err)
			return
		}
		channelName, err := msgCtx.ChannelName()
		if err != nil {
			fmt.Printf("error getting channel name: %v", err)
		}
		channelParams := db.CreateChannelParams{
			ID:          msgCtx.ChannelID,
			ChannelName: channelName,
			WorkspaceID: team.ID,
			CreatedAt:   time.Now(),
			BotToken:    bot.BotToken,
		}

		_, err = dbObj.CreateChannel(context.Background(), channelParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
					fmt.Sprintf("This channel is already set as Stack Overflow Notification channel: *%s*", channelName), false))
				return
				}
			log.Printf("Error creating channel: %v\n", err)
			return
		}

		fmt.Printf("Stack Overflow notification channel is set to %s\n", channelName)
		apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
			fmt.Sprintf("A new Stack Overflow notification channel is set: `%s`", channelName), false))

		// go BotStackOverflow(botCtx, event.ChannelID, "")

		fmt.Printf("A new instance of Stack Overflow channel ID is set to *%s*\n", msgCtx.ChannelID)
	}
}

func AddTagDef(msgCtx *SlackMessageContext, suffix string) {
	apiClient := msgCtx.Api
	event := msgCtx.InnerEvent
	tag := suffix

	if event.Channel != "" && tag != "no_tag" {
		params := db.BindTagParams{
			ChannelID: msgCtx.ChannelID,
			Tag:       tag,
		}

		_, err := dbObj.BindTag(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "violates foreign key constraint") {

				apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
					"Channel is not a notification channel please set it using `soslack_set_so_channel` command", false))
				return
			}
			if strings.Contains(err.Error(), "violates unique constraint") {

				apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
					fmt.Sprintf("This channel is already bound to tag: *%s*", tag), false))
				return
			}
			fmt.Printf("Error binding tag: %v\n", err)
			return
		}

		_, err = dbObj.ActivateTag(context.Background(), tag)
		if err != nil {
			fmt.Printf("Error activating tag: %v\n", err)
			return
		}
		channelName, err := msgCtx.ChannelName()
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		fmt.Printf("Tag *%s* added to %s\n", tag, channelName)
		apiClient.PostMessage(msgCtx.ChannelID, slack.MsgOptionText(
			fmt.Sprintf(
				"New Stack Overflow questions about `%s` will be sent to channel `%s`",
				 tag, channelName), false))

	}
}

func GetUserInfo(msgCtx *SlackMessageContext, suffix string) {
	event := msgCtx.InnerEvent
	apiClient := msgCtx.Api

	userID := event.User
	user, err := apiClient.GetUserInfo(userID)
	if err != nil {
		msgCtx.ReportError(err)
	}
	botID := event.BotID
	bot, err := apiClient.GetBotInfo(botID)
	if err != nil {
		msgCtx.ReportError(err)
	}
	team, err := apiClient.GetTeamInfo()

	if err != nil {
		msgCtx.ReportError(err)
		return
	}
	teamA, err := json.MarshalIndent(team, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		msgCtx.ReportError(err)
		return
	}
	userA, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		msgCtx.ReportError(err)
		return
	}
	botA, err := json.MarshalIndent(bot, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		msgCtx.ReportError(err)
		return
	}

	channelA, err := json.MarshalIndent(event.Channel, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		msgCtx.ReportError(err)
		return
	}
	dataA, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		msgCtx.ReportError(err)
		return
	}

	msgCtx.Reply(fmt.Sprintf("*Team/Workspace*: %s", teamA))
	msgCtx.Reply(fmt.Sprintf("*User*: %s", userA))
	msgCtx.Reply(fmt.Sprintf("*Bot*: %s", botA))
	msgCtx.Reply(fmt.Sprintf("*Channel*: %s", channelA))
	msgCtx.Reply(fmt.Sprintf("*Data*: %s", dataA))
	msgCtx.Reply(fmt.Sprintf("*ChannelID*: %s", msgCtx.ChannelID))
	msgCtx.Reply(fmt.Sprintf("*Type*: %s", event.Type))

}

func ShowTags(msgCtx *SlackMessageContext, suffix string) {
	start := time.Now()
	channelID := msgCtx.ChannelID
	api := msgCtx.Api
	funcStart := time.Now()

	fmt.Printf("ShowTags/GetConversationInfo took %v\n", time.Since(funcStart))
	funcStart = time.Now()

	channelName, err := msgCtx.ChannelName()
	if err != nil {
		fmt.Printf("error getting channel name: %v", err)
		api.PostMessage(channelID, slack.MsgOptionText(
			fmt.Sprintf("an error occured please try again: %v", err), false))
		return
	}

	tags, err := dbObj.GetTagsOfChannel(context.Background(), channelID)
	if err != nil {
		fmt.Printf("Error getting tag of channels: %v\n", err)
		api.PostMessage(channelID, slack.MsgOptionText(
			fmt.Sprintf("Error getting tag of channels: %v\n", err), false))
		return
	}

	fmt.Printf("ShowTags/GetTagsOfChannel took %v\n", time.Since(funcStart))
	funcStart = time.Now()

	if len(tags) == 0 {
		api.PostMessage(channelID, slack.MsgOptionText(
			fmt.Sprintf("No tags bound to channel: %v", channelName), false))
		fmt.Printf("ShowTagsDef took %v\n", time.Since(start))
		return
	}
	tagListStr := ""
	for _, tag := range tags {
		tagListStr += fmt.Sprintf("%v, ", tag.Name)
	}
	api.PostMessage(channelID, slack.MsgOptionText(
		fmt.Sprintf("*Tags bound to channel %v*: %v", channelName, tagListStr), false))

	fmt.Printf("ShowTags/PostMessage took %v\n", time.Since(funcStart))
	fmt.Printf("ShowTagsDef took %v\n", time.Since(start))
}
