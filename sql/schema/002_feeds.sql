-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Reuse existing trigger function
CREATE TRIGGER update_feeds_modtime
BEFORE UPDATE ON feeds
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_feeds_modtime ON feeds;
DROP TABLE IF EXISTS feeds;