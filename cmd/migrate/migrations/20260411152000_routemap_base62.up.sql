CREATE OR REPLACE FUNCTION public.to_base62(input bigint)
RETURNS text
LANGUAGE plpgsql
AS $$
DECLARE
  alphabet text := '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
  base int := 62;
  value bigint := input;
  result text := '';
  remainder int;
BEGIN
  IF value < 0 THEN
    RAISE EXCEPTION 'to_base62 input must be non-negative';
  END IF;

  IF value = 0 THEN
    RETURN '0';
  END IF;

  WHILE value > 0 LOOP
    remainder := (value % base)::int;
    result := substr(alphabet, remainder + 1, 1) || result;
    value := value / base;
  END LOOP;

  RETURN result;
END;
$$;

ALTER TABLE public.routemap
  ALTER COLUMN path DROP NOT NULL;

CREATE OR REPLACE FUNCTION public.set_routemap_path()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.path IS NULL THEN
    NEW.path := public.to_base62(NEW.id);
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS routemap_set_path ON public.routemap;
CREATE TRIGGER routemap_set_path
BEFORE INSERT ON public.routemap
FOR EACH ROW
EXECUTE FUNCTION public.set_routemap_path();
