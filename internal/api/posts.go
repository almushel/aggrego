package api

import (
	"log"
	"net/http"
	"strconv"

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

	ls := r.URL.Query().Get("limit")
	var limit int32
	if len(ls) > 0 {
		l, err := strconv.Atoi(ls)
		if err != nil {
			limit = defaultPostMax
		} else {
			limit = int32(l)
		}
	} else {
		limit = defaultPostMax
	}

	posts, err := api.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  limit,
	})

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
