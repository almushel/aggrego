package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/almushel/aggrego/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *ApiState) PostLikesHandler(w http.ResponseWriter, r *http.Request) {
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
		PostID uuid.UUID `json:"post_id"`
	}
	err = json.Unmarshal(buf, &params)
	if err != nil {
		respondWithError(w, 400, "Malformed request body")
		return
	}

	newLike := database.LikePostParams{
		ID:     uuid.New(),
		PostID: params.PostID,
		UserID: user.ID,
	}

	postLike, err := api.DB.LikePost(r.Context(), newLike)
	if err != nil {
		log.Println("Error: " + err.Error())
		respondWithError(w, 500, "Internal server error"+err.Error())
		return
	}

	respondWithJSON(w, 200, dbToAPI(postLike))
}

func (api *ApiState) DeleteLikesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	likeIDStr := chi.URLParam(r, "postLikeID")
	likeID, _ := uuid.Parse(likeIDStr)
	like, err := api.DB.GetPostLike(r.Context(), likeID)
	if err != nil {
		log.Println("Error:", err)
		respondWithError(w, 404, "Resource not found")
		return
	}

	if like.UserID != user.ID {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	_, err = api.DB.UnlikePost(r.Context(), likeID)
	if err != nil {
		log.Println("Error:", err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, "OK")
}

func (api *ApiState) GetLikesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	likes, err := api.DB.GetUserLikes(r.Context(), user.ID)
	if err != nil {
		log.Println("Error:", err)
		respondWithError(w, 500, "Internal server error")
		return
	}
	var result []PostLike
	for _, like := range likes {
		result = append(result, dbToAPI(like).(PostLike))
	}

	respondWithJSON(w, 200, result)
}
