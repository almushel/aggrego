package api

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/almushel/aggrego/internal/database"
	"github.com/almushel/aggrego/internal/feeds"
	"github.com/google/uuid"
)

const (
	staleFeedstoFetch = 10
	feedFetchInterval = 60 * time.Second
)

func (api *ApiState) StartFetchWorker() {
	for {
		log.Println("Selecting stale feeds")
		dbFeeds, err := api.DB.GetStaleFeeds(context.TODO(), staleFeedstoFetch)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Fetching %d feeds", len(dbFeeds))
			wg := new(sync.WaitGroup)
			for _, feed := range dbFeeds {
				wg.Add(1)

				go func(f database.Feed, w *sync.WaitGroup) {
					defer w.Done()

					rss, err := feeds.FetchRSSFeed(f.Url)
					if err != nil {
						log.Println(err)
					} else {
						for _, post := range rss.Channel.Items {
							_, err := api.DB.CreatePost(context.TODO(), database.CreatePostParams{
								ID:          uuid.New(),
								Title:       post.Title,
								Url:         post.Link,
								Description: post.Description,
								PublishedAt: sql.NullTime{Valid: true, Time: time.Now()},
								FeedID:      f.ID,
							})
							if err != nil && !strings.Contains(err.Error(), "duplicate key value") {
								log.Println(err)
							}
						}
					}
				}(feed, wg)
			}

			wg.Wait()
		}

		time.Sleep(feedFetchInterval)
	}
}
