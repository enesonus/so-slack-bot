-- name: CreateWorkspace :one
INSERT INTO workspaces (id, workspace_name, workspace_domain, bot_id, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;