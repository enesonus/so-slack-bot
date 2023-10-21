-- +goose Up
CREATE TABLE channels (
    id TEXT UNIQUE PRIMARY KEY,
    channel_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    bot_token TEXT NOT NULL REFERENCES bots(bot_token) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE channels;