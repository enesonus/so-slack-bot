// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: channels.sql

package db

import (
	"context"
	"time"
)

const createChannel = `-- name: CreateChannel :one
INSERT INTO channels (id, channel_name, created_at, workspace_id)
VALUES ($1, $2, $3, $4)
RETURNING id, channel_name, created_at, workspace_id, bot_token
`

type CreateChannelParams struct {
	ID          string
	ChannelName string
	CreatedAt   time.Time
	WorkspaceID string
}

func (q *Queries) CreateChannel(ctx context.Context, arg CreateChannelParams) (Channel, error) {
	row := q.db.QueryRowContext(ctx, createChannel,
		arg.ID,
		arg.ChannelName,
		arg.CreatedAt,
		arg.WorkspaceID,
	)
	var i Channel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.CreatedAt,
		&i.WorkspaceID,
		&i.BotToken,
	)
	return i, err
}

const deleteChannel = `-- name: DeleteChannel :one
DELETE FROM channels WHERE id = $1
RETURNING id, channel_name, created_at, workspace_id, bot_token
`

func (q *Queries) DeleteChannel(ctx context.Context, id string) (Channel, error) {
	row := q.db.QueryRowContext(ctx, deleteChannel, id)
	var i Channel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.CreatedAt,
		&i.WorkspaceID,
		&i.BotToken,
	)
	return i, err
}
