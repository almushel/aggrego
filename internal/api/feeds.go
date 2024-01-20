package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/almushel/aggrego/internal/database"
	"github.com/google/uuid"
)

const (
	defaultFeedsPageSize = 20
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

	feedParams := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.URL,
		UserID: user.ID,
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
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	var payload struct {
		Feed Feed       `json:"feed"`
		FF   FeedFollow `json:"feed_follow"`
	}
	payload.Feed = dbToAPI(feed).(Feed)
	payload.FF = dbToAPI(ff).(FeedFollow)

	respondWithJSON(w, 201, payload)
}

func (api *ApiState) GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	var offset int32 = 0
	var limit int32 = defaultFeedsPageSize

	if val, _ := getIntQueryParam(r, "offset"); val >= 0 {
		offset = int32(val)
	}

	if val, _ := getIntQueryParam(r, "limit"); val > 0 {
		limit = int32(val)
	}

	feeds, err := api.DB.GetFeeds(r.Context(), database.GetFeedsParams{
		Offset: offset, Limit: limit,
	})
	if err != nil {
		respondWithError(w, 500, "Internal server error")
	}

	var result []Feed
	for _, f := range feeds {
		result = append(result, dbToAPI(f).(Feed))
	}
	respondWithJSON(w, 200, result)
}
