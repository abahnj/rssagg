-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, feed_id)
);

-- Reuse existing trigger function
CREATE TRIGGER update_feed_follows_modtime
BEFORE UPDATE ON feed_follows
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_feed_follows_modtime ON feed_follows;
DROP TABLE IF EXISTS feed_follows;