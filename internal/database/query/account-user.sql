-- name: GetUserSettings :one
SELECT 
    a.id, a.first_name, a.last_name, a.user_name, a.profile_url, a.email, a.department, a.role, a.is_2fa_enabled, a.is_approved,
    p.account_id, p.units, p.coordinate_format, p.app_theme, p.default_map_view, p.notifications_enabled, p.created_at, p.updated_at
FROM accounts a
LEFT JOIN preferences p ON a.id = p.account_id
WHERE a.id = $1;

-- name: UpsertPreferences :one
INSERT INTO preferences (account_id, units, coordinate_format, app_theme, default_map_view, notifications_enabled, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now())
ON CONFLICT (account_id)
DO UPDATE SET
    units = EXCLUDED.units,
  coordinate_format = EXCLUDED.coordinate_format,
    app_theme = EXCLUDED.app_theme,
    default_map_view = EXCLUDED.default_map_view,
    notifications_enabled = EXCLUDED.notifications_enabled,
    updated_at = now()
RETURNING *;

-- name: ResetPreferences :one
UPDATE preferences
SET
    units = 'metric',
    app_theme = 'light',
coordinate_format = 'dd',
    default_map_view = '{"lat": 6.6745, "lon": -1.5716, "zoom": 12}',
    notifications_enabled = TRUE,
    updated_at = now()
WHERE account_id = $1
RETURNING *;

-- name: UpdateProfile :one
UPDATE accounts
SET
    first_name = $2,
    last_name = $3,
    user_name = $4,
    profile_url = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdatePassword :exec
UPDATE accounts
SET
    password = $2,
    updated_at = now()
WHERE id = $1;

-- name: GetProfile :one
SELECT id, first_name, last_name, user_name, profile_url, email, department, role, is_2fa_enabled, is_approved, password
FROM accounts
WHERE id = $1;
