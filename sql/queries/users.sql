-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
) RETURNING *;


-- name: DeleteAllUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE 
    email = $1;



-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users as u
JOIN refresh_tokens as rf ON u.id = rf.user_id
WHERE token = $1 AND revoked_at IS NULL;