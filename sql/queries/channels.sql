-- name: CreateChannel :one
INSERT INTO channels (id, channel_name, created_at, workspace_id, bot_token)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteChannel :one
DELETE FROM channels WHERE id = $1
RETURNING *;

-- name: GetChannelByBotToken :many
SELECT * FROM channels WHERE bot_token = $1;

-- name: GetChannelByID :one
SELECT * FROM channels WHERE id = $1;