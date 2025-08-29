-- Extensión para bcrypt
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Insertar usuarios de prueba
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES 
  ('admin@site.com', crypt('Admin123!', gen_salt('bf')), 'Admin', 'Root'),
  ('user@site.com', crypt('User123!', gen_salt('bf')), 'Juan', 'Perez')
ON CONFLICT (email) DO NOTHING;