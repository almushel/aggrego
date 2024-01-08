-- +goose Up
ALTER TABLE feeds
ADD last_modified_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DROP last_modified_at;