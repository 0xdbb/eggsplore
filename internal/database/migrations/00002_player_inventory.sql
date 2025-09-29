-- +goose Up
-- +goose StatementBegin

CREATE TABLE players (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id UUID UNIQUE NOT NULL,
  coins BIGINT DEFAULT 0 NOT NULL,
  xp BIGINT DEFAULT 0 NOT NULL,
  level INT DEFAULT 1 NOT NULL,
  settings JSONB DEFAULT '{}'::jsonb, -- preferences (notifications, map theme, etc.)
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);

-- Base inventory table (all items)
CREATE TABLE inventory (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  item_type VARCHAR NOT NULL CHECK (item_type IN ('EGG', 'TOOL', 'BOOST')),
  quantity INT DEFAULT 1 NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Eggs as an inventory subtype
CREATE TABLE eggs (
  inventory_id UUID PRIMARY KEY REFERENCES inventory(id) ON DELETE CASCADE,
  hatched BOOLEAN DEFAULT false,
  type VARCHAR(20) DEFAULT 'BUNNY' NOT NULL, -- e.g., BUNNY, GOLDEN, LEGENDARY
  message TEXT,
  collected_at TIMESTAMPTZ
);

-- Tools as an inventory subtype
CREATE TABLE tools (
  inventory_id UUID PRIMARY KEY REFERENCES inventory(id) ON DELETE CASCADE,
  durability INT NOT NULL DEFAULT 100,
  equipped BOOLEAN DEFAULT false
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS eggs;
DROP TABLE IF EXISTS tools;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
