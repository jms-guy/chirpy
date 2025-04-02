-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;