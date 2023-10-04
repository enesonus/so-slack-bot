-- +goose Up 
CREATE TABLE bots(
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW(),
    bot_token TEXT NOT NULL
);

-- +goose Down
DROP TABLE bots;