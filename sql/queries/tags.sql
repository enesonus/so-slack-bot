-- name: CreateTag :one
INSERT INTO tags (name, has_synonyms, synonyms, is_moderator_only, is_required, count, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ActivateTag :one
UPDATE tags
SET status = 'active'
WHERE name = $1
RETURNING *;

-- name: GetActiveTags :many
SELECT * FROM tags WHERE status = 'active';