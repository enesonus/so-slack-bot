-- +goose Up
CREATE TABLE tag_subscriptions (
    tag TEXT NOT NULL REFERENCES tags(name),
    channel_id TEXT NOT NULL REFERENCES channels(id),
    PRIMARY KEY (tag, channel_id)
);

-- +goose Down
DROP TABLE tag_subscriptions;