-- name: BindTag :one
INSERT INTO tag_subscriptions (tag, channel_id)
VALUES ($1, $2)
RETURNING *;

-- name: UnbindTag :one
DELETE FROM tag_subscriptions
WHERE tag = $1 AND channel_id = $2
RETURNING *;

-- name: GetTagSubscriptions :one
SELECT * FROM tag_subscriptions;

-- name: GetTagSubscriptionsWithName :many
SELECT * FROM tag_subscriptions WHERE tag = $1;

-- name: GetTagSubscriptionsWithChannelId :many
SELECT * FROM tag_subscriptions WHERE channel_id = $1;

-- name: GetSubscriberChannels :one
SELECT 
    channels.*
FROM 
    channels
JOIN 
    tag_subscriptions 
    ON channels.id = tag_subscriptions.channel_id
WHERE 
    tag_subscriptions.tag = $1;