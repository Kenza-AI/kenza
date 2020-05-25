-- +goose Up
CREATE TABLE kenza.accounts (
    id int generated always AS IDENTITY PRIMARY KEY, 
    name varchar NOT NULL CHECK (char_length(name) >= 1),
    email varchar NOT NULL UNIQUE CHECK (char_length(email) >= 3),
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trigger_updated_timestamp
BEFORE UPDATE ON kenza.accounts
FOR EACH ROW
EXECUTE PROCEDURE set_updated_timestamp();

-- +goose Down
DROP TABLE kenza.accounts;