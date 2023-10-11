package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

var pingPong = &slacker.CommandDefinition{
	Description: "Ping!",
	Examples:    []string{"ping", "pang"},
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("pong", slacker.WithThreadReply(true))
	},
}

var echoWordDef = &slacker.CommandDefinition{
	Description: "Echo a word!",
	Examples:    []string{"echo hello"},
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		word := request.Param("word")
		word2 := request.Param("word2")
		response.Reply("word1: " + word + " word2: " + word2)
	},
}

var saySentenceDef = &slacker.CommandDefinition{
	Description: "Say a sentence!",
	Examples:    []string{"say hello there everyone!"},
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		sentence := request.StringParam("sentence", "")
		response.Reply(sentence)
	},
}

var repeatNtimesDef = &slacker.CommandDefinition{
	Description: "Repeat a word a number of times!",
	Examples:    []string{"repeat hello 10"},
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		word := request.StringParam("word", "Hello!")
		number := request.IntegerParam("number", 1)
		for i := 0; i < number; i++ {
			response.Reply(word)
		}
	},
}

var messageReplyDefinition = &slacker.CommandDefinition{
	Description: "Tests errors in new messages",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.ReportError(errors.New("oops, an error occurred"))
	},
}

var uploadDef = &slacker.CommandDefinition{
	Description: "Upload a sentence!",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		sentence := request.Param("sentence")
		apiClient := botCtx.APIClient()
		event := botCtx.Event()

		if event.ChannelID != "" {
			fmt.Printf("Uploading file to channel %s\n", event.ChannelID)
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText("Uploading file ...", false))
			_, err := apiClient.UploadFile(slack.FileUploadParameters{Content: sentence, Channels: []string{event.ChannelID}})
			if err != nil {
				response.ReportError(fmt.Errorf("error encountered when uploading file: %v", err))

			}
		}
	},
}

// ##############################################################################################################################
// bot.Command("upload <sentence>", uploadDef)
// bot.Command("ping", pingPong)
// bot.Command("echo {word} {word2}", echoWordDef)
// bot.Command("say <sentence>", saySentenceDef)
// bot.Command("repeat <word> {number}", repeatNtimesDef)
// bot.Command("message", messageReplyDefinition)

var StackOverflowChannelID = ""
var StackOverflowChannelName = ""

var removeSOChannelDef = &slacker.CommandDefinition{
	Description: "Remove the channel to send SO notifications to",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		apiClient := botCtx.APIClient()
		event := botCtx.Event()

		if event.ChannelID != "" {
			StackOverflowChannelName = ""
			fmt.Printf("Stack Overflow channel is removed\n")
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText("SO question notification channel is removed", false))
			StackOverflowChannelID = ""
			fmt.Printf("Stack Overflow channel ID is removed\n")
		}
	},
}

var setSOChannelDef = &slacker.CommandDefinition{
	Description: "Set the channel to send SO notifications to",
	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		apiClient := botCtx.APIClient()
		event := botCtx.Event()
		tag := request.StringParam("tag", "Hello!")

		if event.ChannelID != "" {
			StackOverflowChannelName = event.Channel.Name
			fmt.Printf("Stack Overflow channel is set to %s\n", StackOverflowChannelName)
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText(
				"SO question notification channel is set to: "+StackOverflowChannelName+" searching for tag: "+tag, false))

			go botStackOverflow(botCtx, event.ChannelID, tag)

			fmt.Printf("A new instance of Stack Overflow channel ID is set to %s\n", event.ChannelID)
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
			StackOverflowChannelName = event.Channel.Name
			fmt.Printf("Tag *%s* added to %s\n", tag, StackOverflowChannelName)
			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText(
				"New Stack Overflow questions about *"+tag+"* will be sent to channel *"+StackOverflowChannelName+"*", false))

			go botStackOverflow(botCtx, event.ChannelID, tag)

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
