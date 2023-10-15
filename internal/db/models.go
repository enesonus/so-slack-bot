// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"database/sql"
	"time"
)

type Bot struct {
	ID             int32
	CreatedAt      time.Time
	LastActivityAt time.Time
	BotToken       string
}

type Channel struct {
	ID          string
	ChannelName string
	BotID       int32
	WorkspaceID string
	CreatedAt   time.Time
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
	TagID int32
	BotID int32
}

type Workspace struct {
	ID              string
	WorkspaceName   sql.NullString
	WorkspaceDomain sql.NullString
	BotID           int32
	CreatedAt       time.Time
}
