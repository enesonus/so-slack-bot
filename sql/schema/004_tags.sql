-- +goose Up
CREATE TABLE tags (
    tag_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    channel_id TEXT NOT NULL,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tags;