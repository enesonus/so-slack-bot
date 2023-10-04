-- name: CreateBot :one
INSERT INTO bots (id, created_at, last_activity_at, bot_token)
VALUES ($1, $2, $3, $4)
RETURNING *;
