-- +goose Up
CREATE TABLE workspaces (
    id TEXT NOT NULL UNIQUE PRIMARY KEY,
    workspace_name VARCHAR(255) NOT NULL,
    workspace_domain VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE workspaces;
