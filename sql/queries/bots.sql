-- name: CreateBot :one
INSERT INTO bots (bot_token, created_at, last_activity_at, workspace_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrCreateBot :one
INSERT INTO bots (bot_token, created_at, last_activity_at, workspace_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (bot_token)
DO UPDATE SET 
    created_at = EXCLUDED.created_at, 
    last_activity_at = EXCLUDED.last_activity_at, 
    workspace_id = EXCLUDED.workspace_id
WHERE 
    bots.created_at != EXCLUDED.created_at OR 
    bots.last_activity_at != EXCLUDED.last_activity_at OR 
    bots.workspace_id != EXCLUDED.workspace_id
RETURNING *;

-- name: GetBotByID :one
SELECT * FROM bots WHERE id = $1;

-- name: GetBotByToken :one
SELECT * FROM bots WHERE bot_token = $1;

-- name: GetBotByWorkspaceID :one
SELECT * FROM bots WHERE workspace_id = $1;

-- name: GetBots :many
SELECT * FROM bots;
