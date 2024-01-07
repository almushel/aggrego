package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/almushel/aggrego/internal/database"
	"github.com/google/uuid"
)

func (api *ApiState) PostFeedsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	buf, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	var params struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	err = json.Unmarshal(buf, &params)
	if err != nil {
		log.Println(err)
		respondWithError(w, 400, "Malformed request body")
		return
	}

	now := time.Now()
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	}

	feed, err := api.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		log.Println(err)
		// NOTE: Is there a better what to handle these?
		if strings.Contains(err.Error(), "duplicate key value") {
			respondWithError(w, 409, "Duplicate feed URL")
		} else {
			respondWithError(w, 500, "Internal server error")
		}
		return
	}

	ff, err := api.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	var payload struct {
		Feed database.Feed       `json:"feed"`
		FF   database.FeedFollow `json:"feed_follow"`
	}
	payload.Feed = feed
	payload.FF = ff

	respondWithJSON(w, 201, payload)
}

func (api *ApiState) GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := api.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, "Internal server error")
	}

	respondWithJSON(w, 200, feeds)
}
