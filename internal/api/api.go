package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

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
