-- +goose Up
-- +goose StatementBegin
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

ALTER TABLE "session" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE-- +goose StatementEnd
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "session" CASCADE
-- +goose StatementEnd
