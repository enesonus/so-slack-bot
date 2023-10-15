-- +goose Up
CREATE TABLE tag_subscriptions (
    tag_id INT NOT NULL REFERENCES tags(id),
    channel_id TEXT NOT NULL REFERENCES channels(id),
    PRIMARY KEY (tag_id, channel_id)
);

-- +goose Down
DROP TABLE tag_subscriptions;