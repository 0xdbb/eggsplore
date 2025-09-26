-- +goose Up
-- +goose StatementBegin
CREATE TABLE "reports" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "location" GEOMETRY(Point, 4326) NOT NULL,
  "locality" varchar,
  "title" varchar,
  description varchar,
  "severity" varchar NOT NULL,
  status varchar NOT NULL DEFAULT 'OPEN',
  "client_ip" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
"updated_at" timestamptz DEFAULT (now())
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "reports" CASCADE
-- +goose StatementEnd
