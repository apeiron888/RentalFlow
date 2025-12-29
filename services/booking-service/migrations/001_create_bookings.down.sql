-- Rollback migration: 001_create_bookings

DROP FUNCTION IF EXISTS generate_booking_number();
DROP TABLE IF EXISTS bookings;
DROP TYPE IF EXISTS cancellation_policy;
DROP TYPE IF EXISTS booking_status;
