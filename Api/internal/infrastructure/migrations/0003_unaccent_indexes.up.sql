-- Habilitar extensión unaccent para búsquedas sin tildes
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Crear envoltorio inmutable sobre unaccent para poder usarlo en índices
-- Nota: referenciamos el diccionario 'unaccent' explícitamente para que sea determinístico
CREATE OR REPLACE FUNCTION public.immutable_unaccent(text)
RETURNS text
LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE
AS $$
    SELECT public.unaccent('public.unaccent'::regdictionary, $1)
$$;

-- Índice funcional para país por nombre sin tildes y en minúsculas
CREATE INDEX IF NOT EXISTS idx_country_unaccent_lower_name
    ON country (public.immutable_unaccent(LOWER(name)));

-- Índice funcional para ciudad por país + nombre sin tildes y en minúsculas
CREATE INDEX IF NOT EXISTS idx_city_country_unaccent_lower_name
    ON city (country_id, public.immutable_unaccent(LOWER(name)));
