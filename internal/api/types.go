package api

import (
	"time"

	"github.com/almushel/aggrego/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Url            string    `json:"url"`
	UserID         uuid.UUID `json:"user_id"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Apikey    string    `json:"apikey"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func dbToAPI(d interface{}) interface{} {
	switch d := d.(type) {
	case database.Feed:
		return Feed{
			d.ID, d.CreatedAt, d.UpdatedAt, d.Name, d.Url, d.UserID, d.LastFetchedAt.Time,
		}
	case database.User:
		return User(d)
	case database.FeedFollow:
		return FeedFollow(d)
	case database.Post:
		return Post{
			d.ID, d.CreatedAt, d.UpdatedAt, d.Title, d.Url, d.Description, d.PublishedAt.Time, d.FeedID,
		}
	}

	return nil
}
