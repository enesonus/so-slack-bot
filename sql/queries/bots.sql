-- name: CreateBot :one
INSERT INTO bots (bot_token, created_at, last_activity_at, workspace_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetBotByID :one
SELECT * FROM bots WHERE id = $1;

-- name: GetBotByToken :one
SELECT * FROM bots WHERE bot_token = $1;

-- name: GetBots :many
SELECT * FROM bots;