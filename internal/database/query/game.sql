-- name: CreateEgg :one
INSERT INTO inventory (player_id, item_type, quantity, description)
VALUES ($1, 'EGG', 1, $2)
RETURNING *;

-- name: AddEggDetails :one
INSERT INTO eggs (inventory_id, type, message, location)
VALUES (
  $1,
  $2,
  $3,
  ST_SetSRID(ST_MakePoint(@lon::float, @lat::float), 4326)
)
RETURNING *;

-- name: GetEggsByPlayer :many
SELECT i.id AS inventory_id, e.type, e.hatched, e.message, e.collected_at
FROM inventory i
JOIN eggs e ON e.inventory_id = i.id
WHERE i.player_id = $1;

-- name: GetToolsByPlayer :many
SELECT i.id AS inventory_id, t.durability, t.equipped, i.description
FROM inventory i
JOIN tools t ON t.inventory_id = i.id
WHERE i.player_id = $1;

-- name: GetInventoryByPlayer :many
SELECT *
FROM inventory
WHERE player_id = $1;

-- name: GetPlayerByAccount :one
SELECT *
FROM players
WHERE account_id = $1;

-- name: UpdatePlayerStats :one
UPDATE players
SET coins = coins + $2,
    xp = xp + $3,
    updated_at = now()
WHERE id = $1
RETURNING *;
