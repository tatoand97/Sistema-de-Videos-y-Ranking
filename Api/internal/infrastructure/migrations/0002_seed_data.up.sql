-- Seed data for initial database population

-- Insert Countries
INSERT INTO country (name, iso_code) VALUES
('Colombia', 'COL'),
('Argentina', 'ARG'),
('México', 'MEX'),
('España', 'ESP'),
('Estados Unidos', 'USA'),
('Brasil', 'BRA'),
('Chile', 'CHL'),
('Perú', 'PER')
ON CONFLICT (iso_code) DO NOTHING;

-- Insert Cities
INSERT INTO city (country_id, name) VALUES
((SELECT country_id FROM country WHERE iso_code = 'COL'), 'Bogotá'),
((SELECT country_id FROM country WHERE iso_code = 'COL'), 'Medellín'),
((SELECT country_id FROM country WHERE iso_code = 'COL'), 'Cali'),
((SELECT country_id FROM country WHERE iso_code = 'ARG'), 'Buenos Aires'),
((SELECT country_id FROM country WHERE iso_code = 'ARG'), 'Córdoba'),
((SELECT country_id FROM country WHERE iso_code = 'MEX'), 'Ciudad de México'),
((SELECT country_id FROM country WHERE iso_code = 'MEX'), 'Guadalajara'),
((SELECT country_id FROM country WHERE iso_code = 'ESP'), 'Madrid'),
((SELECT country_id FROM country WHERE iso_code = 'ESP'), 'Barcelona'),
((SELECT country_id FROM country WHERE iso_code = 'USA'), 'New York'),
((SELECT country_id FROM country WHERE iso_code = 'USA'), 'Los Angeles'),
((SELECT country_id FROM country WHERE iso_code = 'BRA'), 'São Paulo'),
((SELECT country_id FROM country WHERE iso_code = 'BRA'), 'Rio de Janeiro'),
((SELECT country_id FROM country WHERE iso_code = 'CHL'), 'Santiago'),
((SELECT country_id FROM country WHERE iso_code = 'PER'), 'Lima')
ON CONFLICT (country_id, name) DO NOTHING;

-- Insert Roles
INSERT INTO role (name, description) VALUES
('admin', 'Administrador del sistema con acceso completo'),
('moderator', 'Moderador con permisos de gestión de contenido'),
('player', 'Jugador estándar con permisos básicos'),
('viewer', 'Visualizador con permisos de solo lectura')
ON CONFLICT (name) DO NOTHING;

-- Insert Privileges
INSERT INTO privilege (name, description) VALUES
('create_user', 'Crear nuevos usuarios'),
('edit_user', 'Editar información de usuarios'),
('delete_user', 'Eliminar usuarios'),
('view_users', 'Ver lista de usuarios'),
('upload_video', 'Subir videos'),
('edit_video', 'Editar información de videos'),
('delete_video', 'Eliminar videos'),
('view_videos', 'Ver videos'),
('moderate_content', 'Moderar contenido'),
('vote', 'Votar en videos'),
('view_rankings', 'Ver rankings'),
('manage_system', 'Gestionar configuración del sistema')
ON CONFLICT (name) DO NOTHING;

-- Assign Privileges to Roles
INSERT INTO role_privilege (role_id, privilege_id) VALUES
-- Admin: todos los privilegios
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'create_user')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'edit_user')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'delete_user')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'view_users')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'upload_video')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'edit_video')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'delete_video')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'view_videos')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'moderate_content')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'vote')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'view_rankings')),
((SELECT role_id FROM role WHERE name = 'admin'), (SELECT privilege_id FROM privilege WHERE name = 'manage_system')),

-- Moderator: permisos de moderación
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'view_users')),
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'edit_video')),
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'delete_video')),
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'view_videos')),
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'moderate_content')),
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'vote')),
((SELECT role_id FROM role WHERE name = 'moderator'), (SELECT privilege_id FROM privilege WHERE name = 'view_rankings')),

-- Player: permisos de jugador
((SELECT role_id FROM role WHERE name = 'player'), (SELECT privilege_id FROM privilege WHERE name = 'upload_video')),
((SELECT role_id FROM role WHERE name = 'player'), (SELECT privilege_id FROM privilege WHERE name = 'view_videos')),
((SELECT role_id FROM role WHERE name = 'player'), (SELECT privilege_id FROM privilege WHERE name = 'vote')),
((SELECT role_id FROM role WHERE name = 'player'), (SELECT privilege_id FROM privilege WHERE name = 'view_rankings')),

-- Viewer: solo visualización
((SELECT role_id FROM role WHERE name = 'viewer'), (SELECT privilege_id FROM privilege WHERE name = 'view_videos')),
((SELECT role_id FROM role WHERE name = 'viewer'), (SELECT privilege_id FROM privilege WHERE name = 'view_rankings'))
ON CONFLICT (role_id, privilege_id) DO NOTHING;

-- Insert Video Status
INSERT INTO video_status (name, description) VALUES
('trimming', 'Recortando duración a máximo 30 segundos'),
('adjusting_resolution', 'Ajustando resolución y formato de aspecto'),
('adding_watermark', 'Agregando marca de agua ANB'),
('removing_audio', 'Eliminando audio del video'),
('adding_intro_outro', 'Agregando cortinillas de apertura y cierre ANB'),
('processed', 'Video procesado exitosamente, listo para evaluación'),
('failed', 'Error en el procesamiento del video')
ON CONFLICT (name) DO NOTHING;

-- Insert Users (passwords are hashed with bcrypt cost 10)
-- Password for all users: "Password123!"
INSERT INTO users (first_name, last_name, email, password_hash, city_id) VALUES
('Admin', 'Sistema', 'admin@videorank.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Bogotá')),
('Carlos', 'Moderador', 'moderador@videorank.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Madrid')),
('Juan', 'Pérez', 'juan.perez@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Bogotá')),
('María', 'García', 'maria.garcia@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Medellín')),
('Pedro', 'López', 'pedro.lopez@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Buenos Aires')),
('Ana', 'Martínez', 'ana.martinez@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Ciudad de México')),
('Luis', 'Rodríguez', 'luis.rodriguez@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Madrid')),
('Carmen', 'Fernández', 'carmen.fernandez@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'São Paulo')),
('Diego', 'González', 'diego.gonzalez@email.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Santiago')),
('Sofia', 'Viewer', 'viewer@videorank.com', '$2a$10$dbG1mEbaLN9iCH9EmYflK.ddkyJ8aQhw52k5tkbwCzKoEktKOkiQ.', (SELECT city_id FROM city WHERE name = 'Lima'))
ON CONFLICT (email) DO NOTHING;

-- Assign Roles to Users
INSERT INTO user_role (user_id, role_id) VALUES
((SELECT user_id FROM users WHERE email = 'admin@videorank.com'), (SELECT role_id FROM role WHERE name = 'admin')),
((SELECT user_id FROM users WHERE email = 'moderador@videorank.com'), (SELECT role_id FROM role WHERE name = 'moderator')),
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'ana.martinez@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'luis.rodriguez@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'carmen.fernandez@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'diego.gonzalez@email.com'), (SELECT role_id FROM role WHERE name = 'player')),
((SELECT user_id FROM users WHERE email = 'viewer@videorank.com'), (SELECT role_id FROM role WHERE name = 'viewer'))
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Insert Sample Videos
INSERT INTO video (user_id, title, original_file, processed_file, status_id) VALUES
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 'Jugada defensiva destacada', 'juan_video_001.mp4', 'juan_video_001_anb_processed.mp4', (SELECT status_id FROM video_status WHERE name = 'processed')),
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 'Triple desde media cancha', 'juan_video_002.mp4', NULL, (SELECT status_id FROM video_status WHERE name = 'removing_audio')),
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), 'Estrategia ofensiva avanzada', 'maria_video_001.mp4', 'maria_video_001_anb_processed.mp4', (SELECT status_id FROM video_status WHERE name = 'processed')),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), 'Fundamentos de dribleo', 'pedro_video_001.mp4', 'pedro_video_001_anb_processed.mp4', (SELECT status_id FROM video_status WHERE name = 'processed')),
((SELECT user_id FROM users WHERE email = 'ana.martinez@email.com'), 'Competencia nacional ANB', 'ana_video_001.mp4', 'ana_video_001_anb_processed.mp4', (SELECT status_id FROM video_status WHERE name = 'processed')),
((SELECT user_id FROM users WHERE email = 'luis.rodriguez@email.com'), 'Mejores jugadas del torneo', 'luis_video_001.mp4', NULL, (SELECT status_id FROM video_status WHERE name = 'trimming')),
((SELECT user_id FROM users WHERE email = 'carmen.fernandez@email.com'), 'Torneo internacional ANB', 'carmen_video_001.mp4', NULL, (SELECT status_id FROM video_status WHERE name = 'adding_watermark')),
((SELECT user_id FROM users WHERE email = 'diego.gonzalez@email.com'), 'Técnicas profesionales', 'diego_video_001.mp4', NULL, (SELECT status_id FROM video_status WHERE name = 'adding_intro_outro'));

-- Insert Sample Votes
INSERT INTO vote (user_id, video_id) VALUES
-- Votos para el video de Juan (video_id = 1)
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), 1),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), 1),
((SELECT user_id FROM users WHERE email = 'ana.martinez@email.com'), 1),
((SELECT user_id FROM users WHERE email = 'viewer@videorank.com'), 1),

-- Votos para el video de María (video_id = 3)
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 3),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), 3),
((SELECT user_id FROM users WHERE email = 'luis.rodriguez@email.com'), 3),
((SELECT user_id FROM users WHERE email = 'diego.gonzalez@email.com'), 3),
((SELECT user_id FROM users WHERE email = 'viewer@videorank.com'), 3),

-- Votos para el video de Pedro (video_id = 4)
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 4),
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), 4),
((SELECT user_id FROM users WHERE email = 'ana.martinez@email.com'), 4),

-- Votos para el video de Ana (video_id = 5)
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 5),
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), 5),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), 5),
((SELECT user_id FROM users WHERE email = 'luis.rodriguez@email.com'), 5),
((SELECT user_id FROM users WHERE email = 'carmen.fernandez@email.com'), 5),
((SELECT user_id FROM users WHERE email = 'diego.gonzalez@email.com'), 5),
((SELECT user_id FROM users WHERE email = 'viewer@videorank.com'), 5),

-- Votos para el video de Carmen (video_id = 7)
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 7),
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), 7),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), 7),
((SELECT user_id FROM users WHERE email = 'ana.martinez@email.com'), 7),
((SELECT user_id FROM users WHERE email = 'luis.rodriguez@email.com'), 7),

-- Votos para el video de Diego (video_id = 8)
((SELECT user_id FROM users WHERE email = 'juan.perez@email.com'), 8),
((SELECT user_id FROM users WHERE email = 'maria.garcia@email.com'), 8),
((SELECT user_id FROM users WHERE email = 'pedro.lopez@email.com'), 8),
((SELECT user_id FROM users WHERE email = 'ana.martinez@email.com'), 8),
((SELECT user_id FROM users WHERE email = 'luis.rodriguez@email.com'), 8),
((SELECT user_id FROM users WHERE email = 'carmen.fernandez@email.com'), 8)
ON CONFLICT (user_id, video_id) DO NOTHING;