-- +goose Up
-- +goose StatementBegin
ALTER TABLE notifications
ADD COLUMN updated_at TIMESTAMPTZ DEFAULT now() NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notifications
DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
