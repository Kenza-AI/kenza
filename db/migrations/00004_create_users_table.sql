-- +goose Up
CREATE EXTENSION pgcrypto;

CREATE TABLE kenza.users (
    id int generated always AS IDENTITY PRIMARY KEY, 
    username varchar NOT NULL CHECK (char_length(username) >= 1),
    email varchar NOT NULL UNIQUE CHECK (char_length(email) >= 3),
    password varchar NOT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trigger_updated_timestamp
BEFORE UPDATE ON kenza.users
FOR EACH ROW
EXECUTE PROCEDURE set_updated_timestamp();

-- +goose Down
DROP TABLE kenza.users;
DROP EXTENSION pgcrypto;