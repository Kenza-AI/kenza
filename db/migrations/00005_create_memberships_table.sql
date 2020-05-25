-- +goose Up
CREATE TABLE kenza.memberships (
  account_id int NOT NULL REFERENCES kenza.accounts (id) ON UPDATE CASCADE,
  user_id int NOT NULL REFERENCES kenza.users (id) ON UPDATE CASCADE,
  created timestamptz NOT NULL DEFAULT now(), 
  updated timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT memberships_pkey PRIMARY KEY (account_id, user_id)
);

CREATE TRIGGER trigger_updated_timestamp
BEFORE UPDATE ON kenza.memberships
FOR EACH ROW
EXECUTE PROCEDURE set_updated_timestamp();

-- +goose Down
DROP TABLE kenza.memberships;
