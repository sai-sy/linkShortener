DROP POLICY IF EXISTS routemap_delete ON public.routemap;
DROP POLICY IF EXISTS routemap_update ON public.routemap;
DROP POLICY IF EXISTS routemap_insert ON public.routemap;
DROP POLICY IF EXISTS routemap_select ON public.routemap;

ALTER TABLE public.routemap DISABLE ROW LEVEL SECURITY;

DROP FUNCTION IF EXISTS public.authorize(bigint, bigint, text);

DROP INDEX IF EXISTS routemap_workspace_id_idx;

ALTER TABLE public.routemap
  DROP CONSTRAINT IF EXISTS routemap_workspace_id_fkey,
  DROP COLUMN IF EXISTS workspace_id;

DROP TABLE IF EXISTS public.role_permission;
DROP TABLE IF EXISTS public.workspace_member;
DROP TABLE IF EXISTS public.workspace;
