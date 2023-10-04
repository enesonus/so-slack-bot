// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: channels.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createChannel = `-- name: CreateChannel :one
INSERT INTO channels (id, channel_name, bot_id, workspace_id, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, channel_name, bot_id, workspace_id, created_at
`

type CreateChannelParams struct {
	ID          string
	ChannelName string
	BotID       uuid.UUID
	WorkspaceID string
	CreatedAt   time.Time
}

func (q *Queries) CreateChannel(ctx context.Context, arg CreateChannelParams) (Channel, error) {
	row := q.db.QueryRowContext(ctx, createChannel,
		arg.ID,
		arg.ChannelName,
		arg.BotID,
		arg.WorkspaceID,
		arg.CreatedAt,
	)
	var i Channel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.BotID,
		&i.WorkspaceID,
		&i.CreatedAt,
	)
	return i, err
}
