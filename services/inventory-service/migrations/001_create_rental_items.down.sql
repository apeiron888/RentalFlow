-- Rollback migration: 001_create_rental_items

DROP TRIGGER IF EXISTS update_rental_items_updated_at ON rental_items;
DROP TABLE IF EXISTS maintenance_logs;
DROP TABLE IF EXISTS availability_slots;
DROP TABLE IF EXISTS rental_items;
DROP TYPE IF EXISTS maintenance_status;
DROP TYPE IF EXISTS availability_status;
DROP TYPE IF EXISTS item_category;
