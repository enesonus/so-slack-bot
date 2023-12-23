package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func QuestionCheckerAndSender() {
	timePeriod, err := strconv.Atoi(os.Getenv("QUESTION_QUERY_TIME_PERIOD_MINUTES"))
	if err != nil {
		timePeriod = 60
		log.Printf("Couldn't get QUESTION_QUERY_TIME_PERIOD_MINUTES from .env: %v\nDefaulting: %v\n", err.Error(), timePeriod)
	}

	// lastQuestionCheckDate := time.Now().Add(time.Duration(-timePeriod*3) * time.Minute)
	lastQuestionCheckDate := time.Now()

	ticker := time.NewTicker(time.Duration(timePeriod) * time.Minute)
	// ticker := time.NewTicker(time.Duration(5) * time.Second)

	for range ticker.C {
		var questions []StackOverflowQuestion

		if err != nil {
			log.Println("Error getting database: ", err)
		}
		// TODO: active tags might not have any channels subscribed to them
		activeTags, err := dbObj.GetActiveTags(context.Background())
		if err != nil {
			log.Println("Error getting active tags: ", err)
		}

		for _, tag := range activeTags {
			questions = getSOQuestionsAfterTime(tag.Name, lastQuestionCheckDate)
			time.Sleep(2 * time.Second)
			if len(questions) == 0 {
				continue
			}

			tagSubs, err := dbObj.GetTagSubscriptionsWithName(context.Background(), tag.Name)
			if err != nil {
				log.Println("Error getting tag subscriptions: ", err)
			}
			for _, tagSub := range tagSubs {
				channel, err := dbObj.GetChannelByID(context.Background(), tagSub.ChannelID)
				if err != nil {
					log.Println("Error getting channel: ", err)
				}

				for _, question := range questions {

					apiClient := slack.New(channel.BotToken)
					questionTemplate := ">*New %v about %v from %v*!\n"
					decodedName := html.UnescapeString(question.Owner.Display_name)
					markdownName := fmt.Sprintf("<%v|%v>", question.Owner.Link, decodedName)
					markdownQuestion := fmt.Sprintf("<%v|%v>", question.Link, "question")
					questionText := fmt.Sprintf(questionTemplate, markdownQuestion, tag.Name, markdownName)

					apiClient.PostMessage(tagSub.ChannelID, slack.MsgOptionText(questionText, false))
				}

			}
		}

		lastQuestionCheckDate = time.Now()
		fmt.Printf("Last question check date: %v\n", lastQuestionCheckDate)
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
