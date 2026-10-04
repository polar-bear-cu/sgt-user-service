ALTER TABLE users
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS time_in_advanced;
