-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users
SET email = $1,
    hashed_password = $2,
    updated_at = now()
WHERE id = $3
RETURNING id, created_at, updated_at, email;

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET is_chirpy_red = true,
    updated_at = now()
WHERE id = $1;
