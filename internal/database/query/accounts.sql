-- name: GetAccount :one
SELECT * FROM "accounts"
WHERE id = $1;

-- name: GetAccountByEmail :one
SELECT * FROM "accounts"
WHERE LOWER(email) = LOWER($1);


-- name: ListAccounts :many
SELECT *
FROM accounts
WHERE
  ((@role::varchar IS NULL OR @role::varchar = '') OR role::varchar = @role::varchar)
AND ((@status::varchar IS NULL OR @status::varchar = '') OR status::varchar = @status::varchar)
LIMIT $1 OFFSET $2;

-- name: CreateAccount :one
INSERT INTO "accounts" (
    email, password, first_name, last_name, username  
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateAccountRole :one
UPDATE "accounts"
SET role = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdateAccountStatus :one
UPDATE "accounts"
SET status = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM "accounts"
WHERE id = $1;

-- name: UpdateAccountOTP :exec
UPDATE "accounts"
SET otp_code = $2,
    otp_expires_at = $3,
    updated_at = now()
WHERE id = $1;

-- name: ClearAccountOTP :exec
UPDATE "accounts"
SET otp_code = NULL,
    otp_expires_at = NULL,
    updated_at = now()
WHERE id = $1;



-- name: UpdateAccountPassword :one
UPDATE "accounts"
SET password = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountLastActive :exec
UPDATE "accounts"
SET last_active = now()
WHERE id = $1;

-- name: GetAccountMetrics :one
SELECT
    COUNT(*) AS total_accounts,
    SUM(CASE WHEN status = 'ACTIVE' THEN 1 ELSE 0 END) AS active_users,
    SUM(CASE WHEN status = 'INACTIVE' THEN 1 ELSE 0 END) AS inactive_users,
    SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END) AS pending_invites
FROM accounts;


-- name: UpdatePassword :exec
UPDATE accounts
SET
    password = $2,
    updated_at = now()
WHERE id = $1;

-- name: GetProfile :one
SELECT id, first_name, last_name, username, profile_url, email,  role, is_2fa_enabled, is_approved, password
FROM accounts
WHERE id = $1;
