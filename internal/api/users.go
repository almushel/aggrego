package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/almushel/aggrego/internal/database"
	"github.com/google/uuid"
)

func (api *ApiState) UserAuth(w http.ResponseWriter, r *http.Request) (database.User, error) {
	var result database.User
	auth := r.Header.Get("Authorization")
	if len(auth) <= len("ApiKey ") {
		err := errors.New("invalid authorization header")
		respondWithError(w, 401, err.Error())
		return result, err
	}

	apikey := auth[len("ApiKey "):]
	result, err := api.DB.GetUserByKey(r.Context(), apikey)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return result, err
	}

	return result, nil
}

func (api *ApiState) PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	buf, _ := io.ReadAll(r.Body)
	user := struct {
		Name string `json:"name"`
	}{}
	err := json.Unmarshal(buf, &user)
	if err != nil {
		respondWithError(w, 400, "Unexpected request body")
		return
	}

	newUser := database.CreateUserParams{
		ID:   uuid.New(),
		Name: user.Name,
	}

	dbResult, err := api.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 201, dbToAPI(dbResult))
}

func (api *ApiState) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	result, err := api.UserAuth(w, r)
	if err != nil {
		return
	}
	respondWithJSON(w, 200, dbToAPI(result))
}
