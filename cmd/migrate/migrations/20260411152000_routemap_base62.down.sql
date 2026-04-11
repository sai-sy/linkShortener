DROP TRIGGER IF EXISTS routemap_set_path ON public.routemap;
DROP FUNCTION IF EXISTS public.set_routemap_path();

DROP FUNCTION IF EXISTS public.to_base62(bigint);
