-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE "accounts" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "first_name" VARCHAR  CHECK (LENGTH("first_name") > 0),
  "last_name" VARCHAR  CHECK (LENGTH("last_name") > 0),
  "user_name" VARCHAR UNIQUE CHECK (
    LENGTH("user_name") > 0 AND
    "user_name" ~ '^[a-zA-Z0-9_]+$'
  ),
  "email" VARCHAR UNIQUE NOT NULL,
  "department" VARCHAR(255) NOT NULL,
  "password" VARCHAR,
  "status" VARCHAR(20) DEFAULT 'PENDING' NOT NULL,
  "role" VARCHAR NOT NULL,
  "is_2fa_enabled" BOOLEAN DEFAULT TRUE,
  "signup_token" TEXT,
  "signup_token_expires_at" VARCHAR,
  "otp_code" VARCHAR(6),
  "otp_expires_at" VARCHAR,
  "is_approved" BOOLEAN DEFAULT FALSE NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "last_active" TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00 UTC',
  "updated_at" TIMESTAMPTZ DEFAULT now()
);

-- Indexes for performance on frequently fetched fields
CREATE INDEX idx_accounts_email ON accounts (email);
CREATE INDEX idx_accounts_user_name ON accounts (user_name);
CREATE INDEX idx_accounts_role ON accounts (role);
CREATE INDEX idx_accounts_status ON accounts (status);
CREATE INDEX idx_accounts_is_approved ON accounts (is_approved);
CREATE INDEX idx_accounts_department ON accounts (department);
CREATE INDEX idx_accounts_signup_token ON accounts (signup_token);
CREATE INDEX idx_accounts_otp_code ON accounts (otp_code);

CREATE TABLE tasks (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "type" VARCHAR DEFAULT 'PERIODIC' NOT NULL,
  "status" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "aoi_name" VARCHAR(100) NOT NULL,
  "aoi_bbox" GEOMETRY(Polygon, 4326) NOT NULL,
  "image_date" VARCHAR,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_tasks_geom ON tasks USING GIST (aoi_bbox);
CREATE INDEX idx_tasks_status ON tasks (status);

CREATE TABLE mining_segments (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "task_id" UUID NOT NULL,
  "geometry" GEOMETRY(Polygon, 4326) NOT NULL,
  "area" DOUBLE PRECISION NOT NULL,
  "severity" VARCHAR(255) NOT NULL,
  "severity_type" VARCHAR(255) NOT NULL,
  "image_url" TEXT,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_mining_segments_geom ON mining_segments USING GIST (geometry);
CREATE INDEX idx_mining_segments_severity ON mining_segments (severity);
CREATE INDEX idx_mining_segments_type ON mining_segments (severity_type);

CREATE TABLE errors (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "task_id" UUID NOT NULL,
  "error_message" TEXT,
  "timestamp" TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT fk_error_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

COMMENT ON TABLE errors IS 'Log errors for debugging and retry';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS errors;
DROP INDEX IF EXISTS idx_mining_segments_type;
DROP INDEX IF EXISTS idx_mining_segments_severity;
DROP INDEX IF EXISTS idx_mining_segments_geom;
DROP TABLE IF EXISTS mining_segments;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_geom;
DROP TABLE IF EXISTS tasks;
DROP INDEX IF EXISTS idx_accounts_otp_code;
DROP INDEX IF EXISTS idx_accounts_signup_token;
DROP INDEX IF EXISTS idx_accounts_department;
DROP INDEX IF EXISTS idx_accounts_is_approved;
DROP INDEX IF EXISTS idx_accounts_status;
DROP INDEX IF EXISTS idx_accounts_role;
DROP INDEX IF EXISTS idx_accounts_user_name;
DROP INDEX IF EXISTS idx_accounts_email;
DROP TABLE IF EXISTS accounts;
DROP EXTENSION IF EXISTS pgcrypto;
-- +goose StatementEnd
