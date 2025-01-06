-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  fullname,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;
