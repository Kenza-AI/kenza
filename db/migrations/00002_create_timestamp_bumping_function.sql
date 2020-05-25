-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION set_updated_timestamp();