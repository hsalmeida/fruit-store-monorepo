-- Habilita extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Insere o admin somente se não existir
INSERT INTO users (id, username, password_hash, role, created_at, updated_at)
VALUES (
  uuid_generate_v4(),                        -- gera UUID random
  'admin',                                   -- username
  crypt('adminpass', gen_salt('bf')),        -- senha “adminpass” hasheada
  'admin',                                   -- role
  NOW(),                                     -- timestamps
  NOW()
)
ON CONFLICT (username) DO NOTHING;