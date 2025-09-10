----------------------------------------------------------------
-- SEEDS
----------------------------------------------------------------

-- Países
INSERT INTO country (name, iso_code) VALUES
                                         ('Colombia','COL'),('Argentina','ARG'),('México','MEX'),('España','ESP'),
                                         ('Estados Unidos','USA'),('Brasil','BRA'),('Chile','CHL'),('Perú','PER')
    ON CONFLICT (iso_code) DO NOTHING;

-- Ciudades
INSERT INTO city (country_id, name) VALUES
                                        ((SELECT country_id FROM country WHERE iso_code='COL'),'Bogotá'),
                                        ((SELECT country_id FROM country WHERE iso_code='COL'),'Medellín'),
                                        ((SELECT country_id FROM country WHERE iso_code='COL'),'Cali'),
                                        ((SELECT country_id FROM country WHERE iso_code='ARG'),'Buenos Aires'),
                                        ((SELECT country_id FROM country WHERE iso_code='ARG'),'Córdoba'),
                                        ((SELECT country_id FROM country WHERE iso_code='MEX'),'Ciudad de México'),
                                        ((SELECT country_id FROM country WHERE iso_code='MEX'),'Guadalajara'),
                                        ((SELECT country_id FROM country WHERE iso_code='ESP'),'Madrid'),
                                        ((SELECT country_id FROM country WHERE iso_code='ESP'),'Barcelona'),
                                        ((SELECT country_id FROM country WHERE iso_code='USA'),'New York'),
                                        ((SELECT country_id FROM country WHERE iso_code='USA'),'Los Angeles'),
                                        ((SELECT country_id FROM country WHERE iso_code='BRA'),'São Paulo'),
                                        ((SELECT country_id FROM country WHERE iso_code='BRA'),'Rio de Janeiro'),
                                        ((SELECT country_id FROM country WHERE iso_code='CHL'),'Santiago'),
                                        ((SELECT country_id FROM country WHERE iso_code='PER'),'Lima')
    ON CONFLICT (country_id, name) DO NOTHING;

-- Roles
INSERT INTO role (name, description) VALUES
                                         ('admin','Administrador del sistema con acceso completo'),
                                         ('moderator','Moderador con permisos de gestión de contenido'),
                                         ('player','Jugador estándar con permisos básicos'),
                                         ('viewer','Visualizador con permisos de solo lectura')
    ON CONFLICT (name) DO NOTHING;

-- Privilegios
INSERT INTO privilege (name, description) VALUES
                                              ('create_user','Crear nuevos usuarios'),
                                              ('edit_user','Editar información de usuarios'),
                                              ('delete_user','Eliminar usuarios'),
                                              ('view_users','Ver lista de usuarios'),
                                              ('upload_video','Subir videos'),
                                              ('edit_video','Editar información de videos'),
                                              ('delete_video','Eliminar videos'),
                                              ('view_videos','Ver videos'),
                                              ('moderate_content','Moderar contenido'),
                                              ('vote','Votar en videos'),
                                              ('view_rankings','Ver rankings'),
                                              ('manage_system','Gestionar configuración del sistema')
    ON CONFLICT (name) DO NOTHING;

-- Rol-Privilegio
INSERT INTO role_privilege (role_id, privilege_id)
SELECT r.role_id, p.privilege_id
FROM (VALUES
          ('admin','create_user'),('admin','edit_user'),('admin','delete_user'),
          ('admin','view_users'),('admin','upload_video'),('admin','edit_video'),
          ('admin','delete_video'),('admin','view_videos'),('admin','moderate_content'),
          ('admin','vote'),('admin','view_rankings'),('admin','manage_system'),

          ('moderator','view_users'),('moderator','edit_video'),('moderator','delete_video'),
          ('moderator','view_videos'),('moderator','moderate_content'),
          ('moderator','vote'),('moderator','view_rankings'),

          ('player','upload_video'),('player','view_videos'),
          ('player','vote'),('player','view_rankings'),

          ('viewer','view_videos'),('viewer','view_rankings')
     ) AS x(role_name, priv_name)
         JOIN role r ON r.name = x.role_name
         JOIN privilege p ON p.name = x.priv_name
    ON CONFLICT (role_id, privilege_id) DO NOTHING;

-- Users (incluye city_id). Password bcrypt "Password123!"
INSERT INTO users (first_name, last_name, email, password_hash, city_id) VALUES
                                                                             ('Admin','Sistema','admin@videorank.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Bogotá')),
                                                                             ('Carlos','Moderador','moderador@videorank.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Medellín')),
                                                                             ('Juan','Pérez','juan.perez@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Bogotá')),
                                                                             ('María','García','maria.garcia@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Medellín')),
                                                                             ('Pedro','López','pedro.lopez@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Buenos Aires')),
                                                                             ('Ana','Martínez','ana.martinez@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Ciudad de México')),
                                                                             ('Luis','Rodríguez','luis.rodriguez@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Madrid')),
                                                                             ('Carmen','Fernández','carmen.fernandez@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='São Paulo')),
                                                                             ('Diego','González','diego.gonzalez@email.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Santiago')),
                                                                             ('Sofia','Viewer','viewer@videorank.com','$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.',
                                                                              (SELECT city_id FROM city WHERE name='Bogotá'))
    ON CONFLICT (email) DO NOTHING;

-- User-Role
INSERT INTO user_role (user_id, role_id)
SELECT u.user_id, r.role_id
FROM (VALUES
          ('admin@videorank.com','admin'),
          ('moderador@videorank.com','moderator'),
          ('juan.perez@email.com','player'),
          ('maria.garcia@email.com','player'),
          ('pedro.lopez@email.com','player'),
          ('ana.martinez@email.com','player'),
          ('luis.rodriguez@email.com','player'),
          ('carmen.fernandez@email.com','player'),
          ('diego.gonzalez@email.com','player'),
          ('viewer@videorank.com','viewer')
     ) AS x(email, role_name)
         JOIN users u ON u.email = x.email
         JOIN role  r ON r.name  = x.role_name
    ON CONFLICT (user_id, role_id) DO NOTHING;

-- Videos (usar ENUM 'status'; FK por email)
INSERT INTO video (user_id, title, original_file, processed_file, status, processed_at) VALUES
                                                                                            ((SELECT user_id FROM users WHERE email='juan.perez@email.com'),
                                                                                             'Jugada defensiva destacada','juan_video_001.mp4','juan_video_001_anb_processed.mp4','PUBLISHED', now()),
                                                                                            ((SELECT user_id FROM users WHERE email='juan.perez@email.com'),
                                                                                             'Triple desde media cancha','juan_video_002.mp4',NULL,'UPLOADED', NULL),
                                                                                            ((SELECT user_id FROM users WHERE email='maria.garcia@email.com'),
                                                                                             'Estrategia ofensiva avanzada','maria_video_001.mp4','maria_video_001_anb_processed.mp4','PUBLISHED', now()),
                                                                                            ((SELECT user_id FROM users WHERE email='pedro.lopez@email.com'),
                                                                                             'Fundamentos de dribleo','pedro_video_001.mp4',NULL,'UPLOADED', NULL),
                                                                                            ((SELECT user_id FROM users WHERE email='ana.martinez@email.com'),
                                                                                             'Competencia nacional ANB','ana_video_001.mp4','ana_video_001_anb_processed.mp4','PUBLISHED', now()),
                                                                                            ((SELECT user_id FROM users WHERE email='luis.rodriguez@email.com'),
                                                                                             'Mejores jugadas del torneo','luis_video_001.mp4',NULL,'UPLOADED', NULL),
                                                                                            ((SELECT user_id FROM users WHERE email='carmen.fernandez@email.com'),
                                                                                             'Torneo internacional ANB','carmen_video_001.mp4',NULL,'ADDING_WATERMARK', NULL),
                                                                                            ((SELECT user_id FROM users WHERE email='diego.gonzalez@email.com'),
                                                                                             'Técnicas profesionales','diego_video_001.mp4',NULL,'UPLOADED', NULL);

-- Votes (referenciar por original_file para evitar IDs mágicos)
INSERT INTO vote (user_id, video_id)
SELECT u.user_id, v.video_id FROM (
                                      VALUES
                                          ('maria.garcia@email.com','juan_video_001.mp4'),
                                          ('pedro.lopez@email.com','juan_video_001.mp4'),
                                          ('ana.martinez@email.com','juan_video_001.mp4'),
                                          ('viewer@videorank.com','juan_video_001.mp4'),

                                          ('juan.perez@email.com','maria_video_001.mp4'),
                                          ('pedro.lopez@email.com','maria_video_001.mp4'),
                                          ('luis.rodriguez@email.com','maria_video_001.mp4'),
                                          ('diego.gonzalez@email.com','maria_video_001.mp4'),
                                          ('viewer@videorank.com','maria_video_001.mp4'),

                                          ('juan.perez@email.com','pedro_video_001.mp4'),
                                          ('maria.garcia@email.com','pedro_video_001.mp4'),
                                          ('ana.martinez@email.com','pedro_video_001.mp4'),

                                          ('juan.perez@email.com','ana_video_001.mp4'),
                                          ('maria.garcia@email.com','ana_video_001.mp4'),
                                          ('pedro.lopez@email.com','ana_video_001.mp4'),
                                          ('luis.rodriguez@email.com','ana_video_001.mp4'),
                                          ('carmen.fernandez@email.com','ana_video_001.mp4'),
                                          ('diego.gonzalez@email.com','ana_video_001.mp4'),
                                          ('viewer@videorank.com','ana_video_001.mp4'),

                                          ('juan.perez@email.com','carmen_video_001.mp4'),
                                          ('maria.garcia@email.com','carmen_video_001.mp4'),
                                          ('pedro.lopez@email.com','carmen_video_001.mp4'),
                                          ('ana.martinez@email.com','carmen_video_001.mp4'),
                                          ('luis.rodriguez@email.com','carmen_video_001.mp4'),

                                          ('juan.perez@email.com','diego_video_001.mp4'),
                                          ('maria.garcia@email.com','diego_video_001.mp4'),
                                          ('pedro.lopez@email.com','diego_video_001.mp4'),
                                          ('ana.martinez@email.com','diego_video_001.mp4'),
                                          ('luis.rodriguez@email.com','diego_video_001.mp4'),
                                          ('carmen.fernandez@email.com','diego_video_001.mp4')
                                  ) AS x(email, ofile)
                                      JOIN users u ON u.email = x.email
                                      JOIN video v ON v.original_file = x.ofile
    ON CONFLICT (user_id, video_id) DO NOTHING;
