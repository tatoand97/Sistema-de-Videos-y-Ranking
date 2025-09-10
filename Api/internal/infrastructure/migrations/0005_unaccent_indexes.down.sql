-- Eliminar �ndices funcionales y extensi�n unaccent
DROP INDEX IF EXISTS idx_city_country_unaccent_lower_name;
DROP INDEX IF EXISTS idx_country_unaccent_lower_name;

-- Eliminar funci�n envoltorio (si existe)
DROP FUNCTION IF EXISTS public.immutable_unaccent(text);

-- Nota: DROP EXTENSION unaccent; podr�a fallar si otras dependencias existen.
-- Si est�s seguro, descomenta la siguiente l�nea.
-- DROP EXTENSION IF EXISTS unaccent;

