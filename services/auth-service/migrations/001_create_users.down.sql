-- Rollback migration: 001_create_users

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS identity_documents;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS verification_status;
DROP TYPE IF EXISTS user_role;
