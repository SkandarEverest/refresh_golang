-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE(sqlc.narg(email), email),
  username = COALESCE(sqlc.narg(username), username),
  user_image_uri = COALESCE(sqlc.narg(user_image_uri), user_image_uri),
  company_name = COALESCE(sqlc.narg(company_name), company_name),
  company_image_uri = COALESCE(sqlc.narg(company_image_uri), company_image_uri)
WHERE
  id = sqlc.arg(id)
RETURNING *;