-- name: CreateTag :one
INSERT INTO tags (tag_name, created_at, channel_id)
VALUES ($1, $2, $3)
RETURNING *;