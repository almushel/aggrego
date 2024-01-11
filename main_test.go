package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/almushel/aggrego/internal/api"
	"github.com/google/uuid"
)

const (
	apiAddr = "http://localhost:8080/v1"
)

var apikey string
var feedURLs []string
var feedIDs []uuid.UUID
var feedFollowIDs []uuid.UUID

func init() {
	for key, val := range parseEnv() {
		if key == "TESTFEEDS" {
			feedURLs = strings.Split(val, ",")
		}
	}
}

func testRequest(t *testing.T, request *http.Request, codeExpected int, failureMsg string) *http.Response {
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
	for _, url := range feedURLs {
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
	//log.Print("Feeds: ")
	//log.Println(feedIDs)
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
		body = []byte(`{"feed_id":"` + fmt.Sprint(fid) + `"}`)
		request, _ = http.NewRequest("POST", apiAddr+"/feed_follows", bytes.NewBuffer(body))
		request.Header.Add("Authorization", "ApiKey "+apikey)
		response := testRequest(t, request, 201, "Post feed follow failed")

		var ff api.FeedFollow
		if err := unmarshalResponse(response, &ff); err != nil {
			t.Fatal(err)
		}

		feedFollowIDs = append(feedFollowIDs, ff.ID)
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
	request, _ := http.NewRequest("GET", apiAddr+"/posts", nil)
	request.Header.Add("Authorization", "ApiKey "+apikey)
	response := testRequest(t, request, 200, "Failed to get user posts")

	var posts []api.Post
	if err := unmarshalResponse(response, &posts); err != nil {
		t.Fatal("Failed to unmarshel user posts")
	}
}

func TestDeleteFeedFollows(t *testing.T) {
	for _, ffID := range feedFollowIDs {
		request, _ := http.NewRequest("DELETE", apiAddr+"/feed_follows/{"+fmt.Sprint(ffID)+"}", nil)
		request.Header.Add("Authorization", "ApiKey "+apikey)
		testRequest(t, request, 200, "Failed to delete feed follow")
	}
}
