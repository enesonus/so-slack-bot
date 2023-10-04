package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func botInit() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

}

func botStackOverflow(botCtx slacker.BotContext, channelID string, tag string) {
	lastQuestionDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	checkInterval, err := strconv.Atoi(os.Getenv("NEW_QUESTION_CHECK_INTERVAL_SECONDS"))

	if err != nil {
		log.Fatal("Couldn't get checkInterval from .env: " + err.Error())
		checkInterval = 60
	}

	for range time.Tick(time.Duration(checkInterval) * time.Second) {
		var questions []StackOverflowQuestion
		var questionsToPost []StackOverflowQuestion

		if channelID != "" {
			timePeriod, err := strconv.Atoi(os.Getenv("QUESTION_QUERY_TIME_PERIOD_MINUTES"))
			if err != nil {
				log.Println("Couldn't get timePeriod from .env: " + err.Error())
				timePeriod = 60
			}

			fromDate := time.Now().Add(time.Duration(-timePeriod) * time.Minute)
			questions = getSOQuestionsAfterTime(tag, fromDate)
			for _, question := range questions {
				if question.Creation_date > lastQuestionDate.Unix() {
					questionsToPost = append(questionsToPost, question)
				}
			}

			for _, question := range questionsToPost {

				questionTemplate :=
					">*New question about %v from %v*!\nLink: %v \nOwner: %v"
				decodedName := html.UnescapeString(question.Owner.Display_name)
				questionText := fmt.Sprintf(questionTemplate, tag, decodedName, question.Link, question.Owner.Link)

				botCtx.APIClient().PostMessage(channelID, slack.MsgOptionText(questionText, false))

				if question.Creation_date > lastQuestionDate.Unix() {
					lastQuestionDate = time.Unix(question.Creation_date, 0)
				}
			}
			continue
		}

		log.Println("No SO channel set")

	}

}

func printCommandEvents(slackChannel <-chan *slacker.CommandEvent) {
	for event := range slackChannel {
		log.Printf("Command Event Received")
		log.Printf("Command: %v", event.Command)
		log.Printf("Parameters: %v", event.Parameters)
		log.Printf("Event: %v\n\n", event.Event)
	}
}

func getSOQuestionsAfterTime(tag string, fromDate time.Time) []StackOverflowQuestion {
	baseAPIURL := "https://api.stackexchange.com/2.3/questions?key=U4DMV*8nvpm3EOpvf69Rxw(("

	now := time.Now()

	configURL := fmt.Sprintf(
		"&site=stackoverflow&page=1&pagesize=100&fromdate=%v&todate=%v&order=desc&sort=creation",
		fromDate.Unix(), now.Unix())

	tagURL := fmt.Sprintf("&tagged=%v&filter=default", tag)

	url := fmt.Sprintf("%v%v%v", baseAPIURL, configURL, tagURL)

	log.Println("Getting SO questions from: ", url)

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println(err)
		return []StackOverflowQuestion{}
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		log.Println(err)
		return []StackOverflowQuestion{}
	}
	apiResponse := StackExchangeAPIResponse{}

	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return []StackOverflowQuestion{}
	}

	todaysQuestions := []StackOverflowQuestion{}

	todaysQuestions = append(todaysQuestions, apiResponse.Items...)

	return todaysQuestions
}

type StackOverflowQuestion struct {
	Tags  []string
	Owner struct {
		Account_id    int    `json:"account_id"`
		Reputation    int    `json:"reputation"`
		User_id       int    `json:"user_id"`
		User_type     string `json:"user_type"`
		Profile_image string `json:"profile_image"`
		Display_name  string `json:"display_name"`
		Link          string `json:"link"`
	}
	Is_answered        bool   `json:"is_answered"`
	View_count         int    `json:"view_count"`
	Answer_count       int    `json:"answer_count"`
	Score              int    `json:"score"`
	Last_activity_date int64  `json:"last_activity_date"`
	Creation_date      int64  `json:"creation_date"`
	Question_id        int    `json:"question_id"`
	Content_license    string `json:"content_license"`
	Link               string `json:"link"`
	Title              string `json:"title"`
}

type StackExchangeAPIResponse struct {
	Items           []StackOverflowQuestion `json:"items"`
	Has_more        bool                    `json:"has_more"`
	Quota_max       int                     `json:"quota_max"`
	Quota_remaining int                     `json:"quota_remaining"`
}
