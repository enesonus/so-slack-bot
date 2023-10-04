-- +goose Up
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,
    workspace_name TEXT,
    workspace_domain TEXT,
    bot_id UUID NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE workspaces;
