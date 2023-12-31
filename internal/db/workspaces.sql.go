// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: workspaces.sql

package db

import (
	"context"
	"time"
)

const createWorkspace = `-- name: CreateWorkspace :one
INSERT INTO workspaces (id, workspace_name, workspace_domain, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, workspace_name, workspace_domain, created_at
`

type CreateWorkspaceParams struct {
	ID              string
	WorkspaceName   string
	WorkspaceDomain string
	CreatedAt       time.Time
}

func (q *Queries) CreateWorkspace(ctx context.Context, arg CreateWorkspaceParams) (Workspace, error) {
	row := q.db.QueryRowContext(ctx, createWorkspace,
		arg.ID,
		arg.WorkspaceName,
		arg.WorkspaceDomain,
		arg.CreatedAt,
	)
	var i Workspace
	err := row.Scan(
		&i.ID,
		&i.WorkspaceName,
		&i.WorkspaceDomain,
		&i.CreatedAt,
	)
	return i, err
}

const getOrCreateWorkspace = `-- name: GetOrCreateWorkspace :one
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
RETURNING id, workspace_name, workspace_domain, created_at
`

type GetOrCreateWorkspaceParams struct {
	ID              string
	WorkspaceName   string
	WorkspaceDomain string
	CreatedAt       time.Time
}

func (q *Queries) GetOrCreateWorkspace(ctx context.Context, arg GetOrCreateWorkspaceParams) (Workspace, error) {
	row := q.db.QueryRowContext(ctx, getOrCreateWorkspace,
		arg.ID,
		arg.WorkspaceName,
		arg.WorkspaceDomain,
		arg.CreatedAt,
	)
	var i Workspace
	err := row.Scan(
		&i.ID,
		&i.WorkspaceName,
		&i.WorkspaceDomain,
		&i.CreatedAt,
	)
	return i, err
}
