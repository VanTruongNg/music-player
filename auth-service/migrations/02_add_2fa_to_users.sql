-- +goose Up
-- Thêm trường two_fa_enabled và two_fa_secret cho bảng users để hỗ trợ 2FA
ALTER TABLE users
ADD COLUMN two_fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN two_fa_secret VARCHAR(128);

-- +goose Down
ALTER TABLE users
DROP COLUMN two_fa_enabled,
DROP COLUMN two_fa_secret;
