-- +goose Up
CREATE TABLE families (
    id SERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL,
    name VARCHAR(20) NOT NULL UNIQUE,
        CONSTRAINT fk_created_by FOREIGN KEY (created_by)
            REFERENCES users(id),
        CONSTRAINT unique_user_family_name UNIQUE (created_by, name)
);

-- +goose Down
DROP TABLE families;
