-- +goose Up
CREATE TABLE kenza.projects (
    id int generated always AS IDENTITY PRIMARY KEY, 
    account_id int NOT NULL REFERENCES kenza.accounts (id) ON UPDATE CASCADE,
    creator_id int NOT NULL REFERENCES kenza.users (id) ON UPDATE CASCADE,
    title varchar NOT NULL CHECK (char_length(title) >= 1),
    description varchar NOT NULL CHECK (char_length(description) >= 1),
    repository varchar NOT NULL,
    refs varchar NOT NULL DEFAULT '.*',
    vcs_access_token varchar,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trigger_updated_timestamp
BEFORE UPDATE ON kenza.projects
FOR EACH ROW
EXECUTE PROCEDURE set_updated_timestamp();

-- +goose Down
DROP TABLE kenza.projects;