-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: IsUserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = ? LIMIT 1);

-- name: CreateUser :one
INSERT INTO users (id, username, first_name, last_name, created_at) 
VALUES (?, ?, ?, ?, ?) 
RETURNING *;