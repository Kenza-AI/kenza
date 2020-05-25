-- +goose Up
CREATE TABLE kenza.schedules (
    id int generated always AS IDENTITY PRIMARY KEY, 
    project_id int NOT NULL REFERENCES kenza.projects (id) ON UPDATE CASCADE ON DELETE CASCADE,
    title varchar NOT NULL CHECK (char_length(title) >= 1),
    description varchar NOT NULL DEFAULT '',
    ref varchar, -- e.g. refs/heads/master or refs/tags/v1.5.6
    cron_expression varchar,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    UNIQUE (project_id, title)
);

CREATE TRIGGER trigger_updated_timestamp
BEFORE UPDATE ON kenza.schedules
FOR EACH ROW
EXECUTE PROCEDURE set_updated_timestamp();

-- +goose Down
DROP TABLE kenza.schedules;