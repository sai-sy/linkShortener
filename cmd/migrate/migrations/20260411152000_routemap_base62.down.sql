DROP TRIGGER IF EXISTS routemap_set_path ON public.routemap;
DROP FUNCTION IF EXISTS public.set_routemap_path();

ALTER TABLE public.routemap
  ALTER COLUMN path SET NOT NULL;

DROP FUNCTION IF EXISTS public.to_base62(bigint);
