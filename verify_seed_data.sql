-- Script para verificar que los datos seed se cargaron correctamente

-- Verificar países
SELECT 'Countries' as table_name, COUNT(*) as count FROM country;

-- Verificar ciudades
SELECT 'Cities' as table_name, COUNT(*) as count FROM city;

-- Verificar roles
SELECT 'Roles' as table_name, COUNT(*) as count FROM role;

-- Verificar privilegios
SELECT 'Privileges' as table_name, COUNT(*) as count FROM privilege;

-- Verificar asignación de privilegios a roles
SELECT 'Role Privileges' as table_name, COUNT(*) as count FROM role_privilege;

-- Verificar estados de video
SELECT 'Video Status' as table_name, COUNT(*) as count FROM video_status;

-- Verificar usuarios
SELECT 'Users' as table_name, COUNT(*) as count FROM users;

-- Verificar asignación de roles a usuarios
SELECT 'User Roles' as table_name, COUNT(*) as count FROM user_role;

-- Verificar jugadores
SELECT 'Players' as table_name, COUNT(*) as count FROM player;

-- Verificar videos
SELECT 'Videos' as table_name, COUNT(*) as count FROM video;

-- Verificar votos
SELECT 'Votes' as table_name, COUNT(*) as count FROM vote;

-- Mostrar usuarios con sus roles
SELECT 
    u.email,
    u.first_name,
    u.last_name,
    r.name as role_name
FROM users u
JOIN user_role ur ON u.user_id = ur.user_id
JOIN role r ON ur.role_id = r.role_id
ORDER BY u.email;

-- Mostrar videos con información del jugador
SELECT 
    v.title,
    u.first_name || ' ' || u.last_name as player_name,
    c.name as city,
    co.name as country,
    vs.name as status,
    (SELECT COUNT(*) FROM vote vo WHERE vo.video_id = v.video_id) as vote_count
FROM video v
JOIN player p ON v.player_id = p.user_id
JOIN users u ON p.user_id = u.user_id
JOIN city c ON p.city_id = c.city_id
JOIN country co ON c.country_id = co.country_id
JOIN video_status vs ON v.status_id = vs.status_id
ORDER BY vote_count DESC, v.title;