// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"time"
)

type Bot struct {
	ID             int32
	CreatedAt      time.Time
	LastActivityAt time.Time
	BotToken       string
	WorkspaceID    string
}

type Channel struct {
	ID          string
	ChannelName string
	CreatedAt   time.Time
	WorkspaceID string
	BotToken    string
}

type Tag struct {
	ID              int32
	HasSynonyms     bool
	Synonyms        []string
	IsModeratorOnly bool
	IsRequired      bool
	Count           int32
	Name            string
	Status          string
}

type TagSubscription struct {
	Tag       string
	ChannelID string
}

type Workspace struct {
	ID              string
	WorkspaceName   string
	WorkspaceDomain string
	CreatedAt       time.Time
}
