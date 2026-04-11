-- name: GetAllRoutemaps :many
SELECT * FROM public.routemap;

-- name: GetRoutemapsPage :many
SELECT r.id,
       r.path,
       r.destination,
       r.workspace_id,
       w.name AS workspace_name,
       r.created_at
FROM public.routemap r
JOIN public.workspace w ON w.id = r.workspace_id
ORDER BY r.id DESC
LIMIT $1 OFFSET $2;

-- name: InsertRoutemap :exec
INSERT INTO public.routemap (destination, workspace_id)
VALUES ($1, $2);

-- name: GetRoutemap :one
SELECT id, path, destination, workspace_id, created_at
FROM public.routemap
WHERE path = $1
LIMIT 1;

-- name: UpdateRoutemapDestination :exec
UPDATE public.routemap
SET destination = $2
WHERE id = $1;

-- name: CreateAuthUser :one
INSERT INTO auth.users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, password_hash, created_at, updated_at;

-- name: CreateAuthSession :exec
INSERT INTO auth.session (user_id, session_token, csrf_token)
VALUES ($1, $2, $3);

-- name: GetAuthUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM auth.users
WHERE email = $1
LIMIT 1;

-- name: GetAuthSessionByToken :one
SELECT session_token, user_id, csrf_token, created_at
FROM auth.session
WHERE session_token = $1
LIMIT 1;

-- name: CreateProfile :one
INSERT INTO public.profile (user_id, firstname, surname)
VALUES ($1, $2, $3)
RETURNING user_id, firstname, surname, created_at, updated_at;

-- name: CreateWorkspace :one
INSERT INTO public.workspace (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: CreateWorkspaceMember :exec
INSERT INTO public.workspace_member (workspace_id, profile_id, role)
VALUES ($1, $2, $3);

-- name: CreateRolePermission :exec
INSERT INTO public.role_permission (workspace_id, role, permission)
VALUES ($1, $2, $3);

-- name: GetWorkspaceByProfile :one
SELECT w.id, w.name, w.created_at, w.updated_at
FROM public.workspace w
JOIN public.workspace_member wm ON wm.workspace_id = w.id
WHERE wm.profile_id = $1
ORDER BY w.created_at
LIMIT 1;
