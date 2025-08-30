-- Rollback seed data

-- Remove votes
DELETE FROM vote;

-- Remove videos
DELETE FROM video;

-- Remove players
DELETE FROM player;

-- Remove user roles
DELETE FROM user_role;

-- Remove users (except system users if needed)
DELETE FROM users WHERE email NOT LIKE '%@videorank.com' OR email = 'viewer@videorank.com';

-- Remove role privileges
DELETE FROM role_privilege;

-- Remove video status
DELETE FROM video_status;

-- Remove privileges
DELETE FROM privilege;

-- Remove roles
DELETE FROM role;

-- Remove cities
DELETE FROM city;

-- Remove countries
DELETE FROM country;