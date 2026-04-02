-- +goose UP
CREATE TABLE chirp (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "update_at" TIMESTAMP NOT NULL,
    "body" TEXT NOT NULL,
    "user_id" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirp;