package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/enesonus/so-slack-bot/internal/db"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

// ##############################################################################################################################
// bot.Command("upload <sentence>", uploadDef)
// bot.Command("ping", pingPong)
// bot.Command("echo {word} {word2}", echoWordDef)
// bot.Command("say <sentence>", saySentenceDef)
// bot.Command("repeat <word> {number}", repeatNtimesDef)
// bot.Command("message", messageReplyDefinition)

var removeSOChannelDef = &slacker.CommandDefinition{
	Description: "Remove the channel to send SO notifications to",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		apiClient := botCtx.APIClient()
		event := botCtx.Event()

		if event.ChannelID != "" {
			fmt.Printf("Stack Overflow channel is removed\n")
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText("SO question notification channel is removed", false))

			databaseObject, err := db.GetDatabase()
			if err != nil {
				fmt.Printf("Error connecting database: %v\n", err)
			}

			_, err = databaseObject.DeleteChannel(context.Background(), event.ChannelID)
			if err != nil {
				fmt.Printf("Error deleting channel: %v\n", err)
				return
			}
			fmt.Printf("Stack Overflow channel ID is removed\n")
		}
	},
}

var setSOChannelDef = &slacker.CommandDefinition{
	Description: "Set the channel to send SO notifications to",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		apiClient := botCtx.APIClient()
		event := botCtx.Event()

		if event.ChannelID != "" {
			team, err := apiClient.GetTeamInfo()
			if err != nil {
				fmt.Printf("Error getting team info: %v\n", err)
			}
			dbObj, err := db.GetDatabase()
			if err != nil {
				fmt.Printf("Error connecting database: %v\n", err)
			}

			bot, err := dbObj.GetBotByWorkspaceID(context.Background(), team.ID)
			if err != nil {
				fmt.Printf("Error getting bot: %v\n", err)
				return
			}

			channelParams := db.CreateChannelParams{
				ID:          event.ChannelID,
				ChannelName: event.Channel.Name,
				WorkspaceID: team.ID,
				CreatedAt:   time.Now(),
				BotToken:    bot.BotToken,
			}

			_, err = dbObj.CreateChannel(context.Background(), channelParams)
			if err != nil {
				fmt.Printf("Error creating channel: %v\n", err)
				if strings.Contains(err.Error(), "duplicate key value") {
					apiClient.PostMessage(event.ChannelID, slack.MsgOptionText(
						"This channel is already set as Stack Overflow Notification channel: "+event.Channel.Name, false))
				}
				return
			}

			fmt.Printf("Stack Overflow notification channel is set to %s\n", event.Channel.Name)
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText(
				"SO question notification channel is set to: "+event.Channel.Name, false))

			// go BotStackOverflow(botCtx, event.ChannelID, "")

			fmt.Printf("A new instance of Stack Overflow channel ID is set to *%s*\n", event.ChannelID)
		}
	},
}

var addTagDef = &slacker.CommandDefinition{
	Description: "Add Tags to search for in Stack Overflow",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		apiClient := botCtx.APIClient()
		event := botCtx.Event()
		tag := request.StringParam("tag", "no_tag")

		if event.ChannelID != "" && tag != "no_tag" {
			params := db.BindTagParams{
				ChannelID: event.ChannelID,
				Tag:       tag,
			}

			databaseObject, err := db.GetDatabase()
			if err != nil {
				fmt.Printf("Error connecting database: %v\n", err)
			}

			_, err = databaseObject.BindTag(context.Background(), params)
			if err != nil {
				fmt.Printf("Error binding tag: %v\n", err)
				return
			}
			fmt.Printf("Tag *%s* added to %s\n", tag, event.Channel.Name)
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText(
				"New Stack Overflow questions about *"+tag+"* will be sent to channel *"+event.Channel.Name+"*", false))

			// go BotStackOverflow(botCtx, event.ChannelID, tag)

			fmt.Printf("A new instance of Stack Overflow channel ID is set to %s\n", event.ChannelID)
		}
	},
}

var getUserInfoDef = &slacker.CommandDefinition{
	Description: "Get user info",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		userID := botCtx.Event().UserID
		user, err := botCtx.APIClient().GetUserInfo(userID)
		if err != nil {
			response.ReportError(err)
		}
		botID := botCtx.Event().BotID
		bot, err := botCtx.APIClient().GetBotInfo(botID)
		if err != nil {
			response.ReportError(err)
		}
		team, err := botCtx.APIClient().GetTeamInfo()

		if err != nil {
			response.ReportError(err)
			return
		}
		teamA, _ := json.MarshalIndent(team, "", "  ")
		userA, _ := json.MarshalIndent(user, "", "  ")
		botA, _ := json.MarshalIndent(bot, "", "  ")
		channelA, _ := json.MarshalIndent(botCtx.Event().Channel, "", "  ")
		dataA, _ := json.MarshalIndent(botCtx.Event().Data, "", "  ")
		profileA, _ := json.MarshalIndent(botCtx.Event().UserProfile, "", "  ")

		response.Reply(fmt.Sprintf("*Team/Workspace*: %s", teamA))
		response.Reply(fmt.Sprintf("*User*: %s", userA))
		response.Reply(fmt.Sprintf("*Bot*: %s", botA))
		response.Reply(fmt.Sprintf("*Channel*: %s", channelA))
		response.Reply(fmt.Sprintf("*Data*: %s", dataA))
		response.Reply(fmt.Sprintf("*ChannelID*: %s", botCtx.Event().ChannelID))

		response.Reply(fmt.Sprintf("*Type*: %s", botCtx.Event().Type))
		response.Reply(fmt.Sprintf("*UserProfile*: %s", profileA))

	},
}

var showTagsDef = &slacker.CommandDefinition{
	Description: "Get tag info",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		channelID := botCtx.Event().ChannelID
		channelName := botCtx.Event().Channel.Name
		databaseObject, err := db.GetDatabase()

		if err != nil {
			fmt.Printf("Error connecting database: %v\n", err)
			response.Reply(fmt.Sprintf("Error connecting database: %v\n", err))
			return
		}
		tags, err := databaseObject.GetTagsOfChannel(context.Background(), channelID)
		if err != nil {
			fmt.Printf("Error getting channels: %v\n", err)
			response.Reply(fmt.Sprintf("Error getting channels: %v\n", err))
			return
		}
		if len(tags) == 0 {
			response.Reply(fmt.Sprintf("No tags bound to channel: %v", channelName))
			return
		}
		for _, tag := range tags {
			response.Reply(fmt.Sprintf("*Tag bound to %v channel*: %v", channelName, tag.Name))
		}
	},
}
