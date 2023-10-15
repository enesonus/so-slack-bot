-- name: CreateWorkspace :one
INSERT INTO workspaces (id, workspace_name, workspace_domain, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;