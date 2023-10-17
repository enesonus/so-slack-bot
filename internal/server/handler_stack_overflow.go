package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/enesonus/so-slack-bot/internal/bot"
	"github.com/enesonus/so-slack-bot/internal/db"
)

func GetAccessTokenAndStartBot(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	slack_url := "https://slack.com/api/oauth.v2.access"
	client_secret := os.Getenv("SLACK_BOT_CLIENT_SECRET")
	client_id := os.Getenv("SLACK_BOT_CLIENT_ID")
	request_url := fmt.Sprintf("%s?client_secret=%s&client_id=%s&code=%s", slack_url, client_secret, client_id, code)
	res, err := http.Get(request_url)
	if err != nil {
		log.Fatalf("client: could not send request: %s\n", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}

	var resJson AccessTokenAPIResponse
	err = json.Unmarshal(resBody, &resJson)

	if err != nil {
		fmt.Printf("client: could not unmarshal response body: %s\n", err)
	}

	if resJson.OK {
		err = bot.StartSlackBot(resJson.AccessToken)
		if err != nil {
			fmt.Printf("client: could not start slack bot: %s\n", err)
			respondWithJSON(w, 400, map[string]string{"ready": "not ok", "bot_state": "not running", "error": err.Error()})
			return
		}
		respondWithJSON(w, 200, map[string]string{"ready": "ok", "bot_state": "running"})
		return
	}
	var errJson struct {
		OK    string `json:"ok"`
		Error string `json:"error"`
	}
	err = json.Unmarshal(resBody, &errJson)
	if err != nil {
		fmt.Printf("client: could not unmarshal response body: %s\n", err)
	}
	respondWithJSON(w, 400, errJson)

}

func FetchTagsAndSaveToDB(database *db.Queries, startPage int, endPage int) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	hasMore := true
	pageNumber := startPage
	for hasMore && pageNumber <= endPage {
		// Build the request URL
		baseAPIURL := "https://api.stackexchange.com/2.3/tags?key=U4DMV*8nvpm3EOpvf69Rxw(("
		configURL := fmt.Sprintf(
			"&site=stackoverflow&page=%v&pagesize=100&order=asc&sort=popular&filter=!*KkBKr6zsZW(uANU",
			pageNumber)

		url := fmt.Sprintf("%v%v", baseAPIURL, configURL)

		res, err := httpClient.Get(url)
		if err != nil {
			log.Fatal("Error at API request: ", err)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Error reading body of API request: ", err)
		}
		resJson := StackExchangeTagsAPIResponse{}
		err = json.Unmarshal(resBody, &resJson)
		if err != nil {
			log.Fatal("Error unmarshalling body of API request: ", err)
		}
		errCount := 0
		for _, tag := range resJson.Items {

			dbTagParams := db.CreateTagParams{
				Name:            tag.Name,
				HasSynonyms:     tag.HasSynonyms,
				Synonyms:        tag.Synonyms,
				Status:          "passive",
				Count:           int32(tag.Count),
				IsModeratorOnly: tag.IsModeratorOnly,
				IsRequired:      tag.IsRequired,
			}

			_, err = database.CreateTag(context.Background(), dbTagParams)
			if err != nil {
				errCount++
			}
		}
		sleepTime := time.Duration((rand.Intn(20))) * time.Second
		fmt.Print(
			"page: ", pageNumber,
			" hasMore: ", hasMore,
			" errCount: ", errCount,
			" sleepTime: ", sleepTime,
			"\n")
		pageNumber++
		hasMore = resJson.HasMore
		res.Body.Close()
		time.Sleep(sleepTime)
	}
	fmt.Println("Done fetching tags page: ", pageNumber)
}
