-- +goose Up
CREATE TABLE users (
    id BIGINT UNIQUE,
    username VARCHAR,
    firstname VARCHAR,
    joined_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE users;
