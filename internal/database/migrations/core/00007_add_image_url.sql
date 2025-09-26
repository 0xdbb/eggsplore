-- +goose Up
-- +goose StatementBegin
ALTER TABLE accounts
ADD COLUMN profile_url VARCHAR(255) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accounts
DROP COLUMN IF EXISTS profile_url;
-- +goose StatementEnd
