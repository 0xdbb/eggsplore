-- name: GetAccount :one
SELECT * FROM "accounts"
WHERE id = $1;

-- name: GetAccountByEmail :one
SELECT * FROM "accounts"
WHERE LOWER(email) = LOWER($1);

-- name: GetAccountBySignupToken :one
SELECT * FROM "accounts"
WHERE signup_token = $1;

-- name: ListAccounts :many
SELECT *
FROM accounts
WHERE
  ((@role::varchar IS NULL OR @role::varchar = '') OR role::varchar = @role::varchar)
AND ((@status::varchar IS NULL OR @status::varchar = '') OR status::varchar = @status::varchar)
AND ((@department::varchar IS NULL OR @department::varchar = '') OR department::varchar = @department::varchar)
LIMIT $1 OFFSET $2;

-- name: CreateAccount :one
INSERT INTO "accounts" (
    email, first_name, last_name, department,role, status, signup_token, signup_token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, 'PENDING', $6, $7
)
RETURNING *;

-- name: UpdateAccountSetup :one
UPDATE "accounts"
SET first_name = $2,
last_name = $3,
user_name = $4,
    password = $5,
    status = $6,
is_approved = $7,
    signup_token = NULL,
    signup_token_expires_at = NULL,
    updated_at = now(),
    last_active = now()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountRole :one
UPDATE "accounts"
SET role = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountDepartment :one
UPDATE "accounts"
SET department = $2,
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

-- name: UpdateAccountResetToken :exec
UPDATE "accounts"
SET signup_token = $2,
    signup_token_expires_at = $3,
    updated_at = now()
WHERE LOWER( email ) = LOWER( $1 );

-- name: GetAccountByResetToken :one
SELECT id, first_name, last_name, user_name, email, password, role, department, created_at, updated_at, status, signup_token, signup_token_expires_at, otp_code, otp_expires_at
FROM "accounts"
WHERE signup_token = $1;

-- name: UpdateAccountPassword :one
UPDATE "accounts"
SET password = $2,
    signup_token = NULL,
    signup_token_expires_at = NULL,
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
