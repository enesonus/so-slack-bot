-- name: BindTag :one
INSERT INTO tag_subscriptions (tag_id, channel_id)
VALUES ($1, $2)
RETURNING *;