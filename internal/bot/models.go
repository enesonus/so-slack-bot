package bot

import (
	"time"

	"github.com/google/uuid"
)

type Bot struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	LastActivityAt  time.Time `json:"last_activity_at"`
	BotToken        string    `json:"bot_token"`
	WorkspaceID     string    `json:"workspace_id"`
	WorkspaceName   string    `json:"workspace_name"`
	WorkspaceDomain string    `json:"workspace_domain"`
}

type Channel struct {
	ChannelID       string    `json:"channel_id"`
	IssuerID        string    `json:"user_id"`
	BotID           uuid.UUID `json:"bot_id"`
	ChannelName     string    `json:"channel_name"`
	CreatedAt       time.Time `json:"created_at"`
	WorkspaceID     string    `json:"workspace_id"`
	WorkspaceName   string    `json:"workspace_name"`
	WorkspaceDomain string    `json:"workspace_domain"`
	Tags            []string  `json:"tags"`
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
