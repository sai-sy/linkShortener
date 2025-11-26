-- name: GetAllRoutemaps :many
SELECT * FROM public.routemap;

-- name: InsertRoutemap :exec
INSERT INTO public.routemap (path, destination) VALUES ($1, $2);

-- name: GetRoutemap :one
SELECT id, path, destination, created_at
FROM public.routemap
WHERE path = $1
LIMIT 1;

