DROP POLICY IF EXISTS workspace_select ON public.workspace;

ALTER TABLE public.workspace DISABLE ROW LEVEL SECURITY;
