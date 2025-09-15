-- +goose Up
CREATE TABLE user_bank_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    family_id INT NOT NULL,
    token VARCHAR(44) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users(id),
    CONSTRAINT fk_family_id FOREIGN KEY (family_id)
        REFERENCES families(id),
    CONSTRAINT unique_user_token_to_family UNIQUE (user_id, token, family_id)
);

-- +goose Down
DROP TABLE user_bank_tokens;