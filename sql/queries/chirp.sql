-- name: InsertChirp :one
INSERT INTO chirp (id, created_at, update_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
returning *;
