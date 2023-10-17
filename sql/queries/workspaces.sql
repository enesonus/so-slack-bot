-- name: CreateWorkspace :one
INSERT INTO workspaces (id, workspace_name, workspace_domain, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrCreateWorkspace :one
INSERT INTO workspaces (id, workspace_name, workspace_domain, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id)
DO UPDATE SET 
    workspace_name = EXCLUDED.workspace_name, 
    workspace_domain = EXCLUDED.workspace_domain, 
    created_at = EXCLUDED.created_at
WHERE 
    workspaces.workspace_name != EXCLUDED.workspace_name OR 
    workspaces.workspace_domain != EXCLUDED.workspace_domain OR 
    workspaces.created_at != EXCLUDED.created_at
RETURNING *;