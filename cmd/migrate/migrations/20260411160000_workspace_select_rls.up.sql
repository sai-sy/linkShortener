ALTER TABLE public.workspace ENABLE ROW LEVEL SECURITY;

CREATE POLICY workspace_select ON public.workspace
  FOR SELECT
  USING (
    public.authorize(
      current_setting('app.profile_id', true)::bigint,
      id,
      'workspace:read'
    )
  );
