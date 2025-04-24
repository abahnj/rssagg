-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Function needs to be on a single line for Goose to process it correctly
CREATE OR REPLACE FUNCTION update_modified_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_users_modtime ON users;
DROP FUNCTION IF EXISTS update_modified_column();
DROP TABLE IF EXISTS users;