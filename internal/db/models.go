// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Bot struct {
	ID             string
	CreatedAt      time.Time
	LastActivityAt time.Time
	BotToken       string
}

type Channel struct {
	ID          string
	ChannelName string
	BotID       uuid.UUID
	WorkspaceID string
	CreatedAt   time.Time
}

type Tag struct {
	TagName   string
	CreatedAt time.Time
	ChannelID string
}

type Workspace struct {
	ID              string
	WorkspaceName   sql.NullString
	WorkspaceDomain sql.NullString
	BotID           uuid.UUID
	CreatedAt       time.Time
}