-- Habilitar extensi�n unaccent para b�squedas sin tildes
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Crear envoltorio inmutable sobre unaccent para poder usarlo en �ndices
-- Nota: referenciamos el diccionario 'unaccent' expl�citamente para que sea determin�stico
CREATE OR REPLACE FUNCTION public.immutable_unaccent(text)
RETURNS text
LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE
AS $$
    SELECT public.unaccent('public.unaccent'::regdictionary, $1)
$$;

-- �ndice funcional para pa�s por nombre sin tildes y en min�sculas
CREATE INDEX IF NOT EXISTS idx_country_unaccent_lower_name
    ON country (public.immutable_unaccent(LOWER(name)));

-- �ndice funcional para ciudad por pa�s + nombre sin tildes y en min�sculas
CREATE INDEX IF NOT EXISTS idx_city_country_unaccent_lower_name
    ON city (country_id, public.immutable_unaccent(LOWER(name)));

