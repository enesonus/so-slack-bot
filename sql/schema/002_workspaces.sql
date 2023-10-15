-- +goose Up
CREATE TABLE workspaces (
    id TEXT UNIQUE PRIMARY KEY,
    workspace_name VARCHAR(255) NOT NULL,
    workspace_domain VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE workspaces;
