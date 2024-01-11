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
type Post database.Post

func dbToAPI(d interface{}) interface{} {
	switch d.(type) {
	case database.Feed:
		d2, _ := d.(database.Feed)
		result := Feed{
			d2.ID, d2.CreatedAt, d2.UpdatedAt, d2.Name, d2.Url, d2.UserID, d2.LastFetchedAt.Time,
		}
		return result
	case database.User:
		return User(d.(database.User))
	case database.FeedFollow:
		return FeedFollow(d.(database.FeedFollow))
	case database.Post:
		return Post(d.(database.Post))
	}

	return nil
}
