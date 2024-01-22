package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/almushel/aggrego/internal/database"
)

const (
	defaultPostMax = 100
)

func (api *ApiState) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	var offset int32 = 0
	var limit int32 = defaultPostMax

	if val, _ := getIntQueryParam(r, "offset"); val > 0 {
		offset = int32(val)
	}

	if val, _ := getIntQueryParam(r, "limit"); val > 0 {
		limit = int32(val)
	}

	posts, err := api.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Offset: offset,
		Limit:  limit,
	})

	count, err := api.DB.GetPostCount(r.Context(), user.ID)
	if err != nil {
		log.Printf("GetPostCount failed: %s", err)
	} else {
		w.Header().Add("X-Total-Count", fmt.Sprint(count))
	}

	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	var result []Post
	for _, p := range posts {
		result = append(result, dbToAPI(p).(Post))
	}

	respondWithJSON(w, 200, result)
}
