-- name: BindTag :one
INSERT INTO tag_subscriptions (tag, channel_id)
VALUES ($1, $2)
RETURNING *;