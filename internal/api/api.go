package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/almushel/aggrego/internal/database"
)

type ApiState struct {
	DB *database.Queries
}

func NewApi(conn string) (*ApiState, error) {
	result := new(ApiState)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return result, err
	}
	result.DB = database.New(db)

	return result, nil
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	rBody, _ := json.Marshal(payload)
	w.Write(rBody)
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	respondWithJSON(w, status, struct {
		Error string `json:"error"`
	}{Error: msg})
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	type readiness struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, 200, readiness{Status: "ok"})
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal server error")
}

var ErrParamNotFound error = errors.New("query parameter not found")

func getIntQueryParam(r *http.Request, name string) (int, error) {
	result := -1
	if o := r.URL.Query().Get("offset"); len(o) > 0 {
		off, err := strconv.Atoi(o)
		if err != nil {
			return result, err
		} else {
			result = off
		}
	} else {
		return result, ErrParamNotFound
	}

	return result, nil
}
