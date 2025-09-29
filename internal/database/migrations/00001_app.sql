-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE "accounts" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "first_name" VARCHAR  CHECK (LENGTH("first_name") > 0),
  "last_name" VARCHAR  CHECK (LENGTH("last_name") > 0),
  "username" VARCHAR UNIQUE CHECK (
    LENGTH("username") > 0 AND
    "username" ~ '^[a-zA-Z0-9_]+$'
  ),
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "profile_url" VARCHAR(255) DEFAULT NULL,
  "status" VARCHAR(20) DEFAULT 'PENDING' NOT NULL,
  "role" VARCHAR NOT NULL DEFAULT 'USER' ,
  "is_2fa_enabled" BOOLEAN DEFAULT TRUE,
  "otp_code" VARCHAR(6),
  "otp_expires_at" VARCHAR,
  "is_approved" BOOLEAN DEFAULT FALSE NOT NULL,
  "is_verified" BOOLEAN DEFAULT FALSE NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "last_active" TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00 UTC',
  "updated_at" TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE "session" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "account_id" uuid NOT NULL,
  "refresh_token" varchar NOT NULL,
  "account_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz, 
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "session" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

-- Indexes for performance on frequently fetched fields
CREATE INDEX idx_accounts_email ON accounts (email);
CREATE INDEX idx_accounts_username ON accounts (username);
CREATE INDEX idx_accounts_role ON accounts (role);
CREATE INDEX idx_accounts_status ON accounts (status);
CREATE INDEX idx_accounts_is_approved ON accounts (is_approved);
CREATE INDEX idx_accounts_is_verified ON accounts (is_verified);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_accounts_otp_code;
DROP INDEX IF EXISTS idx_accounts_is_approved;
DROP INDEX IF EXISTS idx_accounts_is_verified;
DROP INDEX IF EXISTS idx_accounts_status;
DROP INDEX IF EXISTS idx_accounts_username;
DROP INDEX IF EXISTS idx_accounts_email;
DROP TABLE IF EXISTS "session" CASCADE;
DROP TABLE IF EXISTS accounts;
DROP EXTENSION IF EXISTS pgcrypto;
-- +goose StatementEnd
