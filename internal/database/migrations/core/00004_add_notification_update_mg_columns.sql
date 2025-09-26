-- +goose Up
-- +goose StatementBegin

-- Add new columns to mining_segments
ALTER TABLE mining_segments
ADD COLUMN district TEXT,
ADD COLUMN severity_score INTEGER CHECK (severity_score BETWEEN 1 AND 5),
ADD COLUMN all_violation_types TEXT,
ADD COLUMN last_seen TIMESTAMPTZ DEFAULT now() NOT NULL,
ADD COLUMN proximity_to_water BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN inside_forest_reserve BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN detection_date DATE;

-- Create notifications table
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
  message TEXT NOT NULL,
  type VARCHAR(50) NOT NULL CHECK (
    type IN ('new_site', 'new_report', 'account_approved')
  ),
  related_entity_id UUID,
  related_entity_type VARCHAR(50),
  is_read BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Indexes for notifications
CREATE INDEX idx_notifications_account_id ON notifications (account_id);
CREATE INDEX idx_notifications_is_read ON notifications (is_read);
CREATE INDEX idx_notifications_type ON notifications (type);
CREATE INDEX idx_notifications_related_entity ON notifications (related_entity_type, related_entity_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop notifications table and indexes
DROP INDEX IF EXISTS idx_notifications_related_entity;
DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_account_id;
DROP TABLE IF EXISTS notifications;

-- Remove added columns from mining_segments
ALTER TABLE mining_segments
DROP COLUMN IF EXISTS severity_score,
DROP COLUMN IF EXISTS all_violation_types,
DROP COLUMN IF EXISTS proximity_to_water,
DROP COLUMN IF EXISTS inside_forest_reserve,
DROP COLUMN IF EXISTS detection_date;

-- +goose StatementEnd
