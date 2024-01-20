package api

import (
	"time"

	"github.com/almushel/aggrego/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Url            string
	UserID         uuid.UUID
	LastModifiedAt time.Time
}

type User database.User
type FeedFollow database.FeedFollow

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
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
