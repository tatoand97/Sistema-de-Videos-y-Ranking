BEGIN;

-- 1) Delete seed votes by seed users on seed videos
DELETE FROM vote
WHERE (user_id IN (SELECT user_id FROM users WHERE email IN ('admin@site.com','voter1@site.com'))
       AND video_id IN (
         '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid,
         '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c03'::uuid
       ))
   OR (user_id IN (SELECT user_id FROM users WHERE email='voter1@site.com')
       AND video_id='3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid);

-- 2) Delete seed videos
DELETE FROM video
WHERE video_id IN (
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c01'::uuid,
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c02'::uuid,
  '3f0d4fd2-1d1b-4b0f-9b5e-6a0f2a3b1c03'::uuid
);

-- 3) Delete seed players (or cascade when deleting users)
DELETE FROM player
WHERE user_id IN (
  SELECT user_id FROM users WHERE email IN ('player1@site.com','player2@site.com')
);

-- 4) Delete seed users (cascade removes user_role)
DELETE FROM users
WHERE email IN ('admin@site.com','player1@site.com','player2@site.com','voter1@site.com');

-- 5) Remove privilege assignments from seed roles
DELETE FROM role_privilege
WHERE role_id IN (SELECT role_id FROM role WHERE name IN ('admin','player'));

-- 6) Delete seed roles and privileges
DELETE FROM privilege WHERE name IN ('upload_video','vote_video','manage_users','process_video','view_reports');
DELETE FROM role WHERE name IN ('admin','player');

-- 7) Delete seed cities (if no references)
DELETE FROM city
WHERE name IN ('Bogotá','Medellín','Lima','Cusco','Miami');

-- 8) Delete seed countries (if no references)
DELETE FROM country
WHERE name IN ('Colombia','Peru','United States');

-- 9) Delete seed video statuses (if no references)
DELETE FROM video_status
WHERE name IN ('pending','processing','processed','failed');

COMMIT;
