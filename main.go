package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	. "github.com/almushel/aggrego/internal/api"
	"github.com/almushel/aggrego/internal/util"
)

//go:embed html
var content embed.FS

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	indexPage, err := content.ReadFile("html/index.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
	}

	w.Write(indexPage)
}

func main() {
	var err error
	var corsOptions cors.Options
	var router, v1Router chi.Router
	var api *ApiState
	var server http.Server

	for key, val := range util.ParseEnvFile(".env") {
		os.Setenv(key, val)
	}

	if conn, found := os.LookupEnv("CONN"); found {
		api, err = NewApi(conn)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No CONN value found in .env")
	}

	router = chi.NewRouter()

	corsOptions = cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{""},
		//AllowedMethods: []{},
	}
	router.Use(cors.Handler(corsOptions))
	router.Mount("/", http.HandlerFunc(frontendHandler))

	v1Router = chi.NewRouter()
	v1Router.Get("/readiness", ReadinessHandler)
	v1Router.Get("/error", ErrorHandler)

	v1Router.Post("/users", api.PostUsersHandler)
	v1Router.Get("/users", api.GetUsersHandler)

	v1Router.Post("/feeds", api.PostFeedsHandler)
	v1Router.Get("/feeds", api.GetFeedsHandler)

	v1Router.Post("/feed_follows", api.PostFeedFollowsHandler)
	v1Router.Get("/feed_follows", api.GetFeedFollowsHandler)
	v1Router.Delete("/feed_follows/{feedFollowID}", api.DeleteFeedFollowHandler)
	router.Mount("/v1", v1Router)

	v1Router.Get("/posts", api.GetPostsHandler)

	server.Addr = ":" + os.Getenv("PORT")
	server.Handler = router

	log.Println("Server listening at port " + os.Getenv("PORT"))
	//go api.StartFetchWorker()
	log.Fatal(server.ListenAndServe())
}
