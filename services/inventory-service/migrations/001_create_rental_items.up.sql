-- Inventory Service Database Schema
-- Migration: 001_create_rental_items

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Item categories enum
CREATE TYPE item_category AS ENUM ('vehicle', 'equipment', 'property');

-- Availability status enum
CREATE TYPE availability_status AS ENUM ('available', 'booked', 'maintenance', 'blocked');

-- Maintenance status enum
CREATE TYPE maintenance_status AS ENUM ('scheduled', 'in_progress', 'completed');

-- Rental items table
CREATE TABLE rental_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category item_category NOT NULL,
    subcategory VARCHAR(100),
    
    -- Pricing
    daily_rate DECIMAL(10, 2) NOT NULL,
    weekly_rate DECIMAL(10, 2),
    monthly_rate DECIMAL(10, 2),
    security_deposit DECIMAL(10, 2) DEFAULT 0,
    
    -- Location
    address TEXT,
    city VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Specifications (stored as JSONB for flexibility)
    specifications JSONB DEFAULT '{}',
    
    -- Images (array of URLs)
    images TEXT[] DEFAULT '{}',
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Availability slots table
CREATE TABLE availability_slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rental_item_id UUID NOT NULL REFERENCES rental_items(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status availability_status DEFAULT 'available',
    booking_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Maintenance logs table
CREATE TABLE maintenance_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rental_item_id UUID NOT NULL REFERENCES rental_items(id) ON DELETE CASCADE,
    maintenance_type VARCHAR(50) NOT NULL,  -- scheduled, repair, cleaning
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    cost DECIMAL(10, 2) DEFAULT 0,
    status maintenance_status DEFAULT 'scheduled',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_rental_items_owner ON rental_items(owner_id);
CREATE INDEX idx_rental_items_category ON rental_items(category);
CREATE INDEX idx_rental_items_city ON rental_items(city);
CREATE INDEX idx_rental_items_active ON rental_items(is_active);
CREATE INDEX idx_availability_item ON availability_slots(rental_item_id);
CREATE INDEX idx_availability_dates ON availability_slots(start_date, end_date);
CREATE INDEX idx_availability_status ON availability_slots(status);
CREATE INDEX idx_maintenance_item ON maintenance_logs(rental_item_id);
CREATE INDEX idx_maintenance_status ON maintenance_logs(status);

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for updated_at
CREATE TRIGGER update_rental_items_updated_at
    BEFORE UPDATE ON rental_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
