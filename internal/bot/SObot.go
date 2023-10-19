package bot

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

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func BotStackOverflow(botCtx slacker.BotContext, channelID string, tag string) {
	lastQuestionDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	checkInterval, err := strconv.Atoi(os.Getenv("NEW_QUESTION_CHECK_INTERVAL_SECONDS"))

	if err != nil {
		checkInterval = 60
		log.Printf("Couldn't get NEW_QUESTION_CHECK_INTERVAL_SECONDS from .env: %v\nDefaulting: %v\n", err.Error(), checkInterval)
	}
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)

	for range ticker.C {
		var questions []StackOverflowQuestion
		var questionsToPost []StackOverflowQuestion

		if channelID != "" {
			timePeriod, err := strconv.Atoi(os.Getenv("QUESTION_QUERY_TIME_PERIOD_MINUTES"))
			if err != nil {
				timePeriod = 60
				log.Printf("Couldn't get QUESTION_QUERY_TIME_PERIOD_MINUTES from .env: %v\nDefaulting: %v\n", err.Error(), timePeriod)
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

func PrintCommandEvents(slackChannel <-chan *slacker.CommandEvent) {
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
