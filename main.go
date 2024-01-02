package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/almushel/aggrego/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type apiState struct {
	DB *database.Queries
}

func parseEnv() map[string]string {
	result := make(map[string]string)
	envBuf, err := os.ReadFile(".env")
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(envBuf), "\n") {
		before, after, found := strings.Cut(line, "=")
		if found {
			key := strings.TrimSpace(before)
			val := strings.TrimSpace(after)
			if len(key) > 0 && len(val) > 0 {
				result[key] = val
			}
		}
	}

	return result
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

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	type readiness struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, 200, readiness{Status: "ok"})
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal server error")
}

func (api *apiState) PostUsersHandler(w http.ResponseWriter, r *http.Request) {
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
		ID:        uuid.New(),
		Name:      user.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dbResult, err := api.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, dbResult)
}

func main() {
	var err error
	var corsOptions cors.Options
	var router, v1Router chi.Router
	var api apiState
	var db *sql.DB
	var server http.Server

	for key, val := range parseEnv() {
		os.Setenv(key, val)
	}
	_, found := os.LookupEnv("CONN")
	if !found {
		panic("No CONN value found in .env")
	}
	db, err = sql.Open("postgres", os.Getenv("CONN"))
	if err != nil {
		panic(err)
	}
	api.DB = database.New(db)

	router = chi.NewRouter()

	corsOptions = cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"localhost"},
		//AllowedMethods: []{},
	}
	router.Use(cors.Handler(corsOptions))

	v1Router = chi.NewRouter()
	v1Router.Get("/readiness", readinessHandler)
	v1Router.Get("/error", errorHandler)
	v1Router.Post("/users", api.PostUsersHandler)
	router.Mount("/v1", v1Router)

	server.Addr = "localhost:" + os.Getenv("PORT")
	server.Handler = router

	log.Println("Server listening at port " + os.Getenv("PORT"))
	log.Fatal(server.ListenAndServe())
}
