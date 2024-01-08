package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/almushel/aggrego/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *ApiState) PostFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
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
		FeedID uuid.UUID `json:"feed_id"`
	}
	if err = json.Unmarshal(buf, &params); err != nil {
		respondWithError(w, 400, "Malformed request body")
		return
	}

	now := time.Now()
	ff, err := api.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplciate key") {
			respondWithError(w, 409, "Duplicate feed follow")
		} else {
			respondWithError(w, 500, "Internal server error")
		}
	}

	respondWithJSON(w, 201, dbToAPI(ff))
}

func (api *ApiState) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	ffidStr := chi.URLParam(r, "feedFollowID")
	ffID, err := uuid.Parse(ffidStr)
	if err != nil {
		respondWithError(w, 400, "Invalid feed follow ID")
		return
	}

	err = api.DB.UnfollowFeed(r.Context(), database.UnfollowFeedParams{
		UserID: user.ID,
		ID:     ffID,
	})
	if err != nil {
		respondWithError(w, 409, err.Error())
		return
	}

	respondWithJSON(w, 200, "OK")
}

func (api *ApiState) GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	ff, err := api.DB.GetUserFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, dbToAPI(ff))
}
