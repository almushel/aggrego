package api

import (
	"testing"

	"github.com/almushel/aggrego/internal/database"
)

func TestApitoDB(t *testing.T) {
	var f database.Feed
	var ff database.FeedFollow
	var u database.User
	var p database.Post
	var pl database.PostLike

	t.Run("unsupported type db to api", func(t *testing.T) {
		var arg int
		if result := dbToAPI(arg); result != nil {
			t.Fatalf("Wrong type or nil: %v", result)
		}
	})

	t.Run("feed db to api", func(t *testing.T) {
		if result, ok := dbToAPI(f).(Feed); !ok {
			t.Fatalf("Wrong type or nil: %v", result)
		}
	})

	t.Run("feed follow db to api", func(t *testing.T) {
		if result, ok := dbToAPI(ff).(FeedFollow); !ok {
			t.Fatalf("Wrong type or nil: %v", result)
		}
	})

	t.Run("user db to api", func(t *testing.T) {
		if result, ok := dbToAPI(u).(User); !ok {
			t.Fatalf("Wrong type or nil: %v", result)
		}
	})

	t.Run("post db to api", func(t *testing.T) {
		if result, ok := dbToAPI(p).(Post); !ok {
			t.Fatalf("Wrong type or nil: %v", result)
		}
	})

	t.Run("postLike db to api", func(t *testing.T) {
		if result, ok := dbToAPI(pl).(PostLike); !ok {
			t.Fatalf("Wrong type or nil: %v", result)
		}
	})

}
