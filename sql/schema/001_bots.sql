-- +goose Up 
CREATE TABLE bots (
    id SERIAL UNIQUE PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    bot_token TEXT UNIQUE NOT NULL,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE bots;


