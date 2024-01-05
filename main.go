package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func (api *apiState) UserAuth(w http.ResponseWriter, r *http.Request) (database.User, error) {
	var result database.User
	auth := r.Header.Get("Authorization")
	if len(auth) <= len("ApiKey ") {
		err := errors.New("Invalid authorization header")
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

	now := time.Now()
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      user.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dbResult, err := api.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 201, dbResult)
}

func (api *apiState) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	result, err := api.UserAuth(w, r)
	if err != nil {
		return
	}
	respondWithJSON(w, 200, result)
}

func (api *apiState) PostFeedsHandler(w http.ResponseWriter, r *http.Request) {
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

func (api *apiState) GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := api.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, "Internal server error")
	}

	respondWithJSON(w, 200, feeds)
}

func (api *apiState) PostFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, 201, ff)
}

func (api *apiState) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
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

func (api *apiState) GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.UserAuth(w, r)
	if err != nil {
		return
	}

	ff, err := api.DB.GetUserFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, ff)
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
	v1Router.Get("/users", api.GetUsersHandler)

	v1Router.Post("/feeds", api.PostFeedsHandler)
	v1Router.Get("/feeds", api.GetFeedsHandler)

	v1Router.Post("/feed_follows", api.PostFeedFollowsHandler)
	v1Router.Get("/feed_follows", api.GetFeedFollowsHandler)
	v1Router.Delete("/feed_follows/{feedFollowID}", api.DeleteFeedFollowHandler)
	router.Mount("/v1", v1Router)

	server.Addr = "localhost:" + os.Getenv("PORT")
	server.Handler = router

	log.Println("Server listening at port " + os.Getenv("PORT"))
	log.Fatal(server.ListenAndServe())
}
