-- +goose Up
CREATE TABLE workspaces (
    id TEXT UNIQUE PRIMARY KEY,
    workspace_name TEXT,
    workspace_domain TEXT,
    bot_id TEXT NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE workspaces;
