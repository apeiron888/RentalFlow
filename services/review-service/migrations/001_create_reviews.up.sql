CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE review_type AS ENUM ('renter_to_owner', 'owner_to_renter', 'renter_to_item');

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    target_user_id UUID,
    target_item_id UUID,
    review_type review_type NOT NULL,
    rating DECIMAL(2, 1) CHECK (rating >= 1.0 AND rating <= 5.0),
    comment TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_reviews_booking ON reviews(booking_id);
CREATE INDEX idx_reviews_reviewer ON reviews(reviewer_id);
CREATE INDEX idx_reviews_target_user ON reviews(target_user_id);
CREATE INDEX idx_reviews_target_item ON reviews(target_item_id);
CREATE INDEX idx_reviews_type ON reviews(review_type);
CREATE INDEX idx_reviews_visible ON reviews(is_visible);
