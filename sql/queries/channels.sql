-- name: CreateChannel :one
INSERT INTO channels (id, channel_name, created_at, workspace_id)
VALUES ($1, $2, $3, $4)
RETURNING *;