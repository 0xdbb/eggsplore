-- +goose Up
-- +goose StatementBegin
CREATE TABLE preferences (
    account_id UUID PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
    units VARCHAR(20) DEFAULT 'metric' CHECK (units IN ('metric', 'imperial')),
    app_theme VARCHAR(20) DEFAULT 'light' CHECK (app_theme IN ('light', 'dark', 'system')),
    default_map_view JSONB DEFAULT '{"lat": 6.6745, "lon": -1.5716, "zoom": 12}',
    coordinate_format VARCHAR(20) DEFAULT 'decimal' CHECK (coordinate_format IN ('dd', 'dms')),
    notifications_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE preferences;
-- +goose StatementEnd
