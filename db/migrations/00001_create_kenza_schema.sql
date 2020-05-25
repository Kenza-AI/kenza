-- +goose Up
CREATE SCHEMA IF NOT EXISTS kenza;

-- +goose Down
DROP SCHEMA kenza;