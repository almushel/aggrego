package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func parseEnv() map[string]string {
	result := make(map[string]string)
	envBuf, err := os.ReadFile(".env")
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(envBuf), "\n") {
		key, val, found := strings.Cut(line, "=")
		if found {
			result[key] = val
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

func main() {
	var envParams map[string]string
	var corsOptions cors.Options
	var router, v1Router chi.Router
	var server http.Server

	envParams = parseEnv()
	for key, val := range envParams {
		os.Setenv(key, val)
	}

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
	router.Mount("/v1", v1Router)

	server.Addr = "localhost:" + os.Getenv("PORT")
	server.Handler = router

	log.Println("Server listening at port " + os.Getenv("PORT"))
	log.Fatal(server.ListenAndServe())
}
