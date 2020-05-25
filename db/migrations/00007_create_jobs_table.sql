-- +goose Up
CREATE TABLE kenza.jobs (
    id int generated always AS IDENTITY PRIMARY KEY,
    project_id int NOT NULL REFERENCES kenza.projects (id) ON UPDATE CASCADE ON DELETE CASCADE,
    commit_id varchar NOT NULL,
    delivery_id varchar DEFAULT NULL,
    submitter varchar NOT NULL DEFAULT 'unknown',
    status varchar NOT NULL CHECK (char_length(status) >= 1),
    type varchar NOT NULL DEFAULT 'unknown',
    region varchar DEFAULT NULL,
    endpoint varchar DEFAULT NULL,
    sagemaker_id varchar UNIQUE DEFAULT NULL,
    started timestamptz DEFAULT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trigger_updated_timestamp
BEFORE UPDATE ON kenza.jobs
FOR EACH ROW
EXECUTE PROCEDURE set_updated_timestamp();

-- +goose Down
DROP TABLE kenza.jobs;
