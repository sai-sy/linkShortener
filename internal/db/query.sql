-- name: GetAllRoutemaps :many
SELECT * FROM public.routemap;

-- name: InsertRoutemap :exec
INSERT INTO public.routemap (path, destination) VALUES ($1, $2);

-- name: GetRoutemap :one
SELECT id, path, destination, created_at
FROM public.routemap
WHERE path = $1
LIMIT 1;

-- name: CreateAuthUser :one
INSERT INTO auth.users (id, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, password_hash, created_at, updated_at;

-- name: GetAuthUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM auth.users
WHERE email = $1
LIMIT 1;

-- name: CreateProfile :one
INSERT INTO public.profile (user_id, firstname, surname)
VALUES ($1, $2, $3)
RETURNING user_id, firstname, surname, created_at, updated_at;
