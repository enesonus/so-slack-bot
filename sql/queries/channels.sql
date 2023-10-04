-- name: CreateChannel :one
INSERT INTO channels (id, channel_name, bot_id, workspace_id, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
