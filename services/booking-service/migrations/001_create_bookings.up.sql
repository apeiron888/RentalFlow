-- Booking Service Database Schema
-- Migration: 001_create_bookings

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Booking status enum
CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'active', 'completed', 'cancelled');

-- Cancellation policy enum
CREATE TYPE cancellation_policy AS ENUM ('flexible', 'moderate', 'strict');

-- Bookings table
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_number VARCHAR(20) UNIQUE NOT NULL,
    renter_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    rental_item_id UUID NOT NULL,
    status booking_status DEFAULT 'pending',
    
    -- Dates
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    total_days INTEGER NOT NULL,
    
    -- Pricing
    daily_rate DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    security_deposit DECIMAL(10, 2) DEFAULT 0,
    service_fee DECIMAL(10, 2) DEFAULT 0,
    total_amount DECIMAL(10, 2) NOT NULL,
    
    -- Locations
    pickup_address TEXT,
    pickup_notes TEXT,
    pickup_time TIMESTAMP WITH TIME ZONE,
    return_address TEXT,
    return_notes TEXT,
    return_time TIMESTAMP WITH TIME ZONE,
    
    -- Additional services (stored as JSONB)
    additional_services JSONB DEFAULT '[]',
    
    cancellation_policy cancellation_policy DEFAULT 'moderate',
    
    -- Agreement
    agreement_signed BOOLEAN DEFAULT FALSE,
    agreement_url TEXT,
    renter_signature BYTEA,
    owner_signature BYTEA,
    
    -- Cancellation
    cancelled_by UUID,
    cancellation_reason TEXT,
    
    -- Payment
    payment_status VARCHAR(50),
    payment_id UUID,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_bookings_renter ON bookings(renter_id);
CREATE INDEX idx_bookings_owner ON bookings(owner_id);
CREATE INDEX idx_bookings_rental_item ON bookings(rental_item_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_dates ON bookings(start_date, end_date);
CREATE INDEX idx_bookings_number ON bookings(booking_number);

-- Function to generate booking number
CREATE OR REPLACE FUNCTION generate_booking_number() RETURNS TEXT AS $$
DECLARE
    new_number TEXT;
BEGIN
    new_number := 'BK' || TO_CHAR(NOW(), 'YYYYMMDD') || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0');
    RETURN new_number;
END;
$$ LANGUAGE plpgsql;
