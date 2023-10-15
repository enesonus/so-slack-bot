package server

type AuthedUser struct {
	ID string `json:"id"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccessTokenAPIResponse struct {
	OK                  bool       `json:"ok"`
	AppID               string     `json:"app_id"`
	AuthedUser          AuthedUser `json:"authed_user"`
	Scope               string     `json:"scope"`
	TokenType           string     `json:"token_type"`
	AccessToken         string     `json:"access_token"`
	BotUserID           string     `json:"bot_user_id"`
	Team                Team       `json:"team"`
	Enterprise          *string    `json:"enterprise"`
	IsEnterpriseInstall bool       `json:"is_enterprise_install"`
}

type TagItem struct {
	HasSynonyms     bool     `json:"has_synonyms"`
	Synonyms        []string `json:"synonyms"`
	IsModeratorOnly bool     `json:"is_moderator_only"`
	IsRequired      bool     `json:"is_required"`
	Count           int      `json:"count"`
	Name            string   `json:"name"`
}

type StackExchangeTagsAPIResponse struct {
	Items          []TagItem `json:"items"`
	HasMore        bool      `json:"has_more"`
	QuotaMax       int       `json:"quota_max"`
	QuotaRemaining int       `json:"quota_remaining"`
}
