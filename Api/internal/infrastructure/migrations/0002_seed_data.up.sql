BEGIN;

-- Ensure extension for UUID and bcrypt
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-----------------------------
--  Country
-----------------------------
INSERT INTO country (name, iso_code) VALUES
  ('Colombia','COL'),
  ('Peru','PER'),
  ('United States','USA')
ON CONFLICT (name) DO NOTHING;

-----------------------------
--  City  (use NOT EXISTS for idempotency)
-----------------------------
-- Bogotá (Colombia)
INSERT INTO city (country_id, name)
SELECT p.country_id, 'Bogotá'
FROM country p WHERE p.name='Colombia'
AND NOT EXISTS (
  SELECT 1 FROM city c
  WHERE c.name='Bogotá' AND c.country_id=p.country_id
);

-- Medellín (Colombia)
INSERT INTO city (country_id, name)
SELECT p.country_id, 'Medellín'
FROM country p WHERE p.name='Colombia'
AND NOT EXISTS (
  SELECT 1 FROM city c
  WHERE c.name='Medellín' AND c.country_id=p.country_id
);

-- Lima (Peru)
INSERT INTO city (country_id, name)
SELECT p.country_id, 'Lima'
FROM country p WHERE p.name='Peru'
AND NOT EXISTS (
  SELECT 1 FROM city c
  WHERE c.name='Lima' AND c.country_id=p.country_id
);

-- Cusco (Peru)
INSERT INTO city (country_id, name)
SELECT p.country_id, 'Cusco'
FROM country p WHERE p.name='Peru'
AND NOT EXISTS (
  SELECT 1 FROM city c
  WHERE c.name='Cusco' AND c.country_id=p.country_id
);

-- Miami (United States)
INSERT INTO city (country_id, name)
SELECT p.country_id, 'Miami'
FROM country p WHERE p.name='United States'
AND NOT EXISTS (
  SELECT 1 FROM city c
  WHERE c.name='Miami' AND c.country_id=p.country_id
);

-----------------------------
--  Role
-----------------------------
INSERT INTO role (name, description) VALUES
  ('admin','System administrator'),
  ('player','Participating user who uploads and votes on videos')
ON CONFLICT (name) DO NOTHING;

-----------------------------
--  Privilege
-----------------------------
INSERT INTO privilege (name, description) VALUES
  ('upload_video','Allows uploading videos'),
  ('vote_video','Allows voting on videos'),
  ('manage_users','Allows managing users'),
  ('process_video','Allows operating the processing pipeline'),
  ('view_reports','Allows viewing reports')
ON CONFLICT (name) DO NOTHING;

-----------------------------
--  Video_Status
-----------------------------
INSERT INTO video_status (name, description) VALUES
  ('pending','Pending processing'),
  ('processing','Currently processing'),
  ('processed','Successfully processed'),
  ('failed','Processing failed')
ON CONFLICT (name) DO NOTHING;

-----------------------------
--  Users  (bcrypt via pgcrypto: crypt('pwd', gen_salt('bf')))
--  Assumes columns: email, password_hash, first_name, last_name
-----------------------------
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES 
  ('admin@site.com',    crypt('Admin123!',    gen_salt('bf')), 'Admin',  'Root'),
  ('player1@site.com',  crypt('Secreta123',   gen_salt('bf')), 'Juan',   'Perez'),
  ('player2@site.com',  crypt('Secreta123',   gen_salt('bf')), 'Maria',  'Gomez'),
  ('voter1@site.com',   crypt('Secreta123',   gen_salt('bf')), 'Carlos', 'Lopez')
ON CONFLICT (email) DO NOTHING;

-----------------------------
--  Player (1:1 Users)
-----------------------------
-- player1 -> Bogotá
INSERT INTO player (user_id, city_id)
SELECT u.user_id, c.city_id
FROM users u
JOIN city c ON c.name='Bogotá'
WHERE u.email='player1@site.com'
AND NOT EXISTS (
  SELECT 1 FROM player j WHERE j.user_id = u.user_id
);

-- player2 -> Lima
INSERT INTO player (user_id, city_id)
SELECT u.user_id, c.city_id
FROM users u
JOIN city c ON c.name='Lima'
WHERE u.email='player2@site.com'
AND NOT EXISTS (
  SELECT 1 FROM player j WHERE j.user_id = u.user_id
);

-----------------------------
--  User_Role (M:N)
-----------------------------
-- admin -> admin
INSERT INTO user_role (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, role r
WHERE u.email='admin@site.com' AND r.name='admin'
ON CONFLICT (user_id, role_id) DO NOTHING;

-- player1 -> player
INSERT INTO user_role (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, role r
WHERE u.email='player1@site.com' AND r.name='player'
ON CONFLICT (user_id, role_id) DO NOTHING;

-- player2 -> player
INSERT INTO user_role (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, role r
WHERE u.email='player2@site.com' AND r.name='player'
ON CONFLICT (user_id, role_id) DO NOTHING;

-- voter1 -> (optional) player to allow voting
INSERT INTO user_role (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, role r
WHERE u.email='voter1@site.com' AND r.name='player'
ON CONFLICT (user_id, role_id) DO NOTHING;

-----------------------------
--  Role_Privilege (M:N)
-----------------------------
-- admin: all privileges
INSERT INTO role_privilege (role_id, privilege_id)
SELECT r.role_id, p.privilege_id
FROM role r CROSS JOIN privilege p
WHERE r.name='admin'
ON CONFLICT (role_id, privilege_id) DO NOTHING;

-- player: upload + vote
INSERT INTO role_privilege (role_id, privilege_id)
SELECT r.role_id, p.privilege_id
FROM role r
JOIN privilege p ON p.name IN ('upload_video','vote_video')
WHERE r.name='player'
ON CONFLICT (role_id, privilege_id) DO NOTHING;

-----------------------------
--  Video (fixed UUIDs for idempotency)
-----------------------------
-- Video 1 (player1, processed)
INSERT INTO video (video_id, player_id, title, original_file, processed_file, status_id, uploaded_at, processed_at)
SELECT
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid,
  j.user_id,
  'Long-range screamer',
  's3://bucket/videos/p1_screamer.mp4',
  's3://bucket/videos/p1_screamer_proc.mp4',
  e.status_id,
  now(),
  now()
FROM player j
JOIN users u ON u.user_id=j.user_id AND u.email='player1@site.com'
JOIN video_status e ON e.name='processed'
WHERE NOT EXISTS (
  SELECT 1 FROM video v WHERE v.video_id='3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid
);

-- Video 2 (player1, pending)
INSERT INTO video (video_id, player_id, title, original_file, processed_file, status_id, uploaded_at)
SELECT
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c02'::uuid,
  j.user_id,
  'Brilliant play and assist',
  's3://bucket/videos/p1_assist.mp4',
  NULL,
  e.status_id,
  now()
FROM player j
JOIN users u ON u.user_id=j.user_id AND u.email='player1@site.com'
JOIN video_status e ON e.name='pending'
WHERE NOT EXISTS (
  SELECT 1 FROM video v WHERE v.video_id='3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c02'::uuid
);

-- Video 3 (player2, processing)
INSERT INTO video (video_id, player_id, title, original_file, processed_file, status_id, uploaded_at)
SELECT
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c03'::uuid,
  j.user_id,
  'Epic save',
  's3://bucket/videos/p2_save.mp4',
  NULL,
  e.status_id,
  now()
FROM player j
JOIN users u ON u.user_id=j.user_id AND u.email='player2@site.com'
JOIN video_status e ON e.name='processing'
WHERE NOT EXISTS (
  SELECT 1 FROM video v WHERE v.video_id='3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c03'::uuid
);

-----------------------------
--  Vote (uses unique constraint for idempotency)
-----------------------------
-- admin votes Video 1 and 3
INSERT INTO vote (user_id, video_id)
SELECT ua.user_id, v.video_id
FROM users ua, video v
WHERE ua.email='admin@site.com' AND v.video_id IN (
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid,
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c03'::uuid
)
ON CONFLICT ON CONSTRAINT unique_vote_user_video DO NOTHING;

-- voter1 votes Video 1
INSERT INTO vote (user_id, video_id)
SELECT ua.user_id, v.video_id
FROM users ua, video v
WHERE ua.email='voter1@site.com' AND v.video_id='3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid
ON CONFLICT ON CONSTRAINT unique_vote_user_video DO NOTHING;

COMMIT;
