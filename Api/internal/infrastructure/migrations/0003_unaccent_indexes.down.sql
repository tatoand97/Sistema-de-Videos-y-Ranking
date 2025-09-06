-- Eliminar índices funcionales y extensión unaccent
DROP INDEX IF EXISTS idx_city_country_unaccent_lower_name;
DROP INDEX IF EXISTS idx_country_unaccent_lower_name;

-- Eliminar función envoltorio (si existe)
DROP FUNCTION IF EXISTS public.immutable_unaccent(text);

-- Nota: DROP EXTENSION unaccent; podría fallar si otras dependencias existen.
-- Si estás seguro, descomenta la siguiente línea.
-- DROP EXTENSION IF EXISTS unaccent;
