-- +goose Up
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    has_synonyms BOOLEAN NOT NULL,
    synonyms VARCHAR(255)[],
    is_moderator_only BOOLEAN NOT NULL,
    is_required BOOLEAN NOT NULL,
    count INT NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE tags;