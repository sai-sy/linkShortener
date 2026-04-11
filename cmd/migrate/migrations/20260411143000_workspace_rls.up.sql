CREATE TABLE IF NOT EXISTS public.workspace (
  id uuid DEFAULT uuidv7() PRIMARY KEY,
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.workspace_member (
  workspace_id uuid NOT NULL REFERENCES public.workspace(id) ON DELETE CASCADE,
  profile_id uuid NOT NULL REFERENCES public.profile(user_id) ON DELETE CASCADE,
  role text NOT NULL CHECK (role IN ('owner', 'admin', 'member')),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (workspace_id, profile_id)
);

CREATE TABLE IF NOT EXISTS public.role_permission (
  workspace_id uuid NOT NULL REFERENCES public.workspace(id) ON DELETE CASCADE,
  role text NOT NULL,
  permission text NOT NULL,
  UNIQUE (workspace_id, role, permission)
);

ALTER TABLE public.routemap
  ADD COLUMN IF NOT EXISTS workspace_id uuid;

INSERT INTO public.workspace (id, name)
SELECT p.user_id,
       COALESCE(NULLIF(p.firstname, ''), u.email) || '''s Workspace'
FROM public.profile p
JOIN auth.users u ON u.id = p.user_id
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.workspace_member (workspace_id, profile_id, role)
SELECT p.user_id, p.user_id, 'owner'
FROM public.profile p
ON CONFLICT DO NOTHING;

INSERT INTO public.role_permission (workspace_id, role, permission)
SELECT w.id, perms.role, perms.permission
FROM public.workspace w
CROSS JOIN (
  VALUES
    ('owner', 'workspace:read'),
    ('owner', 'workspace:update'),
    ('owner', 'routemap:create'),
    ('owner', 'routemap:read'),
    ('owner', 'routemap:update'),
    ('owner', 'routemap:delete'),
    ('member', 'workspace:read'),
    ('member', 'routemap:read')
) AS perms(role, permission)
ON CONFLICT DO NOTHING;

UPDATE public.routemap
SET workspace_id = (
  SELECT id
  FROM public.workspace
  ORDER BY created_at
  LIMIT 1
)
WHERE workspace_id IS NULL;

ALTER TABLE public.routemap
  ALTER COLUMN workspace_id SET NOT NULL,
  ADD CONSTRAINT routemap_workspace_id_fkey
    FOREIGN KEY (workspace_id) REFERENCES public.workspace(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS routemap_workspace_id_idx
  ON public.routemap(workspace_id);

CREATE OR REPLACE FUNCTION public.authorize(
  profile_id uuid,
  workspace_id_param uuid,
  permission text
) RETURNS boolean
LANGUAGE sql
AS $$
  SELECT profile_id IS NOT NULL
  AND EXISTS (
    SELECT 1
    FROM public.workspace_member wm
    JOIN public.role_permission rp
      ON rp.workspace_id = wm.workspace_id
     AND rp.role = wm.role
    WHERE wm.profile_id = profile_id
      AND wm.workspace_id = workspace_id_param
      AND rp.permission = permission
  );
$$;

ALTER TABLE public.routemap ENABLE ROW LEVEL SECURITY;

CREATE POLICY routemap_select ON public.routemap
  FOR SELECT
  USING (
    public.authorize(
      current_setting('app.profile_id', true)::uuid,
      workspace_id,
      'routemap:read'
    )
  );

CREATE POLICY routemap_insert ON public.routemap
  FOR INSERT
  WITH CHECK (
    public.authorize(
      current_setting('app.profile_id', true)::uuid,
      workspace_id,
      'routemap:create'
    )
  );

CREATE POLICY routemap_update ON public.routemap
  FOR UPDATE
  USING (
    public.authorize(
      current_setting('app.profile_id', true)::uuid,
      workspace_id,
      'routemap:update'
    )
  )
  WITH CHECK (
    public.authorize(
      current_setting('app.profile_id', true)::uuid,
      workspace_id,
      'routemap:update'
    )
  );

CREATE POLICY routemap_delete ON public.routemap
  FOR DELETE
  USING (
    public.authorize(
      current_setting('app.profile_id', true)::uuid,
      workspace_id,
      'routemap:delete'
    )
  );
