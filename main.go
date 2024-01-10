package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	. "github.com/almushel/aggrego/internal/api"
	"github.com/almushel/aggrego/internal/database"
	"github.com/almushel/aggrego/internal/feeds"
)

const (
	staleFeedstoFetch = 10
)

func parseEnv() map[string]string {
	result := make(map[string]string)
	envBuf, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
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

func startFetchWorker(api *ApiState) {
	for {
		log.Println("Selecting stale feeds")
		dbFeeds, err := api.DB.GetStaleFeeds(context.TODO(), staleFeedstoFetch)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Fetching %d feeds", len(dbFeeds))
			wg := new(sync.WaitGroup)
			for _, feed := range dbFeeds {
				go func(f database.Feed, w *sync.WaitGroup) {
					w.Add(1)
					defer w.Done()

					rss, err := feeds.FetchRSSFeed(f.Url)
					if err != nil {
						log.Println(err)
					} else {
						for _, post := range rss.Channel.Items {
							println(post.Title)
						}
					}
				}(feed, wg)
			}

			wg.Wait()
		}

		time.Sleep(60 * time.Second)
	}
}

func main() {
	var err error
	var corsOptions cors.Options
	var router, v1Router chi.Router
	var api *ApiState
	var server http.Server

	for key, val := range parseEnv() {
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
		AllowedOrigins:   []string{"localhost"},
		//AllowedMethods: []{},
	}
	router.Use(cors.Handler(corsOptions))

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

	server.Addr = "localhost:" + os.Getenv("PORT")
	server.Handler = router

	log.Println("Server listening at port " + os.Getenv("PORT"))
	go startFetchWorker(api)
	log.Fatal(server.ListenAndServe())
}
