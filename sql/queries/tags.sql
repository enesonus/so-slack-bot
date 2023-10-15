-- name: CreateTag :one
INSERT INTO tags (name, has_synonyms, synonyms, is_moderator_only, is_required, count, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
