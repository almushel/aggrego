package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/almushel/aggrego/internal/api"
	"github.com/almushel/aggrego/internal/util"
)

var apiAddr string
var apikey string
var feedURLs []string
var feedIDs []uuid.UUID
var feedFollowIDs []uuid.UUID
var postIDs []uuid.UUID
var postLikeIDs []uuid.UUID

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		for key, val := range util.ParseEnvFile(".env") {
			if key == "TESTFEEDS" {
				feedURLs = strings.Split(val, ",")
			} else if key == "PORT" {
				apiAddr = "http://:" + val + "/v1"
			}
		}
	}
	if apiAddr == "" {
		apiAddr = "http://localhost:8080"
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func testRequest(t *testing.T, request *http.Request, codeExpected int, failureMsg string) *http.Response {
	if testing.Short() {
		t.Skip()
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	} else if response.StatusCode != codeExpected {
		t.Fatal(failureMsg + ": " + response.Status)
	}

	return response
}

func unmarshalResponse(r *http.Response, body interface{}) error {
	buf, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, body)

	return err
}

func TestPostUser(t *testing.T) {
	body := []byte(`{"name": "testUser"}`)
	request, err := http.NewRequest("POST", apiAddr+"/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	response := testRequest(t, request, 201, "User post failed")

	var user api.User
	if err = unmarshalResponse(response, &user); err != nil {
		t.Fatal(err)
	}

	apikey = user.Apikey
}

func TestGetUser(t *testing.T) {
	request, _ := http.NewRequest("GET", apiAddr+"/users", nil)

	testRequest(t, request, 401, "Get user succeeded without api key")

	request.Header.Add("Authorization", "ApiKey "+apikey)
	response := testRequest(t, request, 200, "Get user failed")

	var user api.User
	if err := unmarshalResponse(response, &user); err != nil {
		t.Fatal(err)
	}
}

func TestPostFeed(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	for _, url := range feedURLs {
		t.Run("Post "+url, func(t *testing.T) {
			body := []byte(`{"name":"testfeed", "url":"` + url + `"}`)
			request, _ := http.NewRequest("POST", apiAddr+"/feeds", bytes.NewBuffer(body))
			request.Header.Add("Authorization", "ApiKey "+apikey)

			response := testRequest(t, request, 201, "Failed to post feed")

			var payload struct {
				Feed api.Feed       `json:"feed"`
				FF   api.FeedFollow `json:"feed_follow"`
			}
			if err := unmarshalResponse(response, &payload); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetFeeds(t *testing.T) {
	request, _ := http.NewRequest("GET", apiAddr+"/feeds", nil)
	response := testRequest(t, request, 200, "Failed to get feeds")

	var feeds []api.Feed
	if err := unmarshalResponse(response, &feeds); err != nil {
		t.Fatal(err)
	}

	for _, feed := range feeds {
		feedIDs = append(feedIDs, feed.ID)
	}
}

func TestPostFeedFollow(t *testing.T) {
	body := []byte(`{"name": "testUser2"}`)
	request, _ := http.NewRequest("POST", apiAddr+"/users", bytes.NewBuffer(body))
	response := testRequest(t, request, 201, "Failed to post testUser2")

	var user api.User
	if err := unmarshalResponse(response, &user); err != nil {
		t.Fatal(err)
	}
	apikey = user.Apikey

	for _, fid := range feedIDs {
		fidStr := fmt.Sprint(fid)
		t.Run("Follow feed "+fidStr, func(t *testing.T) {
			body = []byte(`{"feed_id":"` + fidStr + `"}`)
			request, _ = http.NewRequest("POST", apiAddr+"/feed_follows", bytes.NewBuffer(body))
			request.Header.Add("Authorization", "ApiKey "+apikey)
			response := testRequest(t, request, 201, "Post feed follow failed")

			var ff api.FeedFollow
			if err := unmarshalResponse(response, &ff); err != nil {
				t.Fatal(err)
			}

			feedFollowIDs = append(feedFollowIDs, ff.ID)
		})
	}
}

func TestGetFeedFollows(t *testing.T) {
	request, _ := http.NewRequest("GET", apiAddr+"/feed_follows", nil)
	request.Header.Add("Authorization", "ApiKey "+apikey)
	response := testRequest(t, request, 200, "Failed to get feed follows")

	var ff []api.FeedFollow
	if err := unmarshalResponse(response, &ff); err != nil {
		t.Fatal(err)
	}
}

func TestGetPosts(t *testing.T) {
	var offset int = 0
	var limit int = 20

	for {
		var pageLength int
		name := fmt.Sprintf("Get posts %d-%d", offset, offset+limit)
		t.Run(name, func(t *testing.T) {
			requestURL := fmt.Sprintf("%s/posts?offset=%d&limit=%d", apiAddr, offset, limit)
			request, _ := http.NewRequest("GET", requestURL, nil)
			request.Header.Add("Authorization", "ApiKey "+apikey)
			response := testRequest(t, request, 200, "Failed to get user posts")

			var posts []api.Post
			if err := unmarshalResponse(response, &posts); err != nil {
				t.Fatal("Failed to unmarshel user posts")
			}
			pageLength = len(posts)

			for _, p := range posts {
				postIDs = append(postIDs, p.ID)
			}
		})

		if pageLength > limit {
			t.Fatalf("Page length wanted: %d, received: %d", limit, pageLength)
		} else if pageLength < limit {
			break
		}

		offset += pageLength
	}

}

func TestDeleteFeedFollows(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	for _, ffID := range feedFollowIDs {
		ffidStr := fmt.Sprint(ffID)
		t.Run("Unfollow feed "+ffidStr, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", apiAddr+"/feed_follows/{"+ffidStr+"}", nil)
			request.Header.Add("Authorization", "ApiKey "+apikey)
			testRequest(t, request, 200, "Failed to delete feed follow")
		})
	}
}

func TestPostLikes(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	for _, pid := range postIDs {
		pidStr := pid.String()
		t.Run("Like post "+pidStr, func(t *testing.T) {
			body := []byte(`{"post_id":"` + pidStr + `"}`)
			request, _ := http.NewRequest("POST", apiAddr+"/post_likes", bytes.NewBuffer(body))
			request.Header.Add("Authorization", "ApiKey "+apikey)
			response := testRequest(t, request, 200, "Failed to POST post like "+pidStr)

			var newLike api.PostLike
			err := unmarshalResponse(response, newLike)
			if err == nil {
				postLikeIDs = append(postLikeIDs, newLike.ID)
			}
		})
	}
}

func TestGetLikes(t *testing.T) {
	request, _ := http.NewRequest("GET", apiAddr+"/post_likes", nil)
	request.Header.Add("Authorization", "ApiKey "+apikey)
	testRequest(t, request, 200, "Failed to get liked posts for "+apikey)
}

func TestDeleteLikes(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	for _, pid := range postLikeIDs {
		pidStr := pid.String()
		t.Run("Unlike post "+pidStr, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", apiAddr+"/post_likes/"+pidStr, nil)
			request.Header.Add("Authorization", "ApiKey "+apikey)
			testRequest(t, request, 200, "Failed to DELETE post like "+pidStr)
		})
	}
}
