package feeds

import (
	"strings"
	"testing"

	"github.com/almushel/aggrego/internal/util"
)

func TestFetchFeed(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var feeds []string
	env := util.ParseEnvFile("../../.env")

	if testFeeds, ok := env["TESTFEEDS"]; ok {
		feeds = strings.Split(testFeeds, ",")
	} else {
		t.Fatal("No testfeeds in .env")
	}

	if len(feeds) == 0 {
		t.Fatal("No feeds to retrieve")
	}

	for _, url := range feeds {
		if feed, err := FetchRSSFeed(url); err != nil {
			t.Fatal(err)
		} else {
			for _, post := range feed.Channel.Items {
				println(post.Title)
			}
		}
	}
}
