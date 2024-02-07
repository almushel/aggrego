-- +goose Up
CREATE TABLE post_likes (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	post_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY(post_id) REFERENCES posts(id),
	UNIQUE(user_id, post_id)
);

-- +goose Down
DROP TABLE post_likes;