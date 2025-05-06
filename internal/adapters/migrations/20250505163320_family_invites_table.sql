-- +goose Up
CREATE TABLE family_invites (
    id SERIAL PRIMARY KEY,
    family_id INT NOT NULL,
    code VARCHAR(6) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_family_id FOREIGN KEY (family_id)
        REFERENCES families(id),
    CONSTRAINT fk_created_by FOREIGN KEY (created_by)
        REFERENCES users(id)
);

-- +goose Down
DROP TABLE family_invites;
