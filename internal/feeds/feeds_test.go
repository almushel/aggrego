package feeds

import (
	"os"
	"strings"
	"testing"
)

func TestFeed(t *testing.T) {
	var feeds []string
	if buf, err := os.ReadFile("testfeeds.env"); err != nil {
		t.Fatal(err)
	} else {
		for _, feed := range strings.Split(string(buf), "\n") {
			feeds = append(feeds, feed)
		}
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
