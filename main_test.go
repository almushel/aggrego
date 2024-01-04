package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/almushel/aggrego/internal/database"
	"github.com/google/uuid"
)

const (
	apiAddr = "http://localhost:8080/v1"
)

var apikey string

func TestPostUser(t *testing.T) {
	var rBody struct {
		Name string `json:"name"`
	}
	rBody.Name = "testUser"
	rBuf, err := json.Marshal(rBody)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", apiAddr+"/users", bytes.NewBuffer(rBuf))
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 201 {
		t.Fatal(response.Status)
	}

	result, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var user database.User
	err = json.Unmarshal(result, &user)
	if err != nil {
		t.Fatal(err)
	}

	apikey = user.Apikey
}

func TestGetUser(t *testing.T) {
	request, _ := http.NewRequest("GET", apiAddr+"/users", nil)

	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != 401 {
		t.Fatal("Get user succeeded without api key")
	}

	request.Header.Add("Authorization", "ApiKey "+apikey)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal(response.Status)
	}

	buf, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	var user database.User
	err = json.Unmarshal(buf, &user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostFeed(t *testing.T) {
	body := []byte(`{"name": "testfeed", "url": "http://test.com/` + fmt.Sprint(uuid.New()) + `"}`)
	request, _ := http.NewRequest("POST", apiAddr+"/feeds", bytes.NewBuffer(body))
	request.Header.Add("Authorization", "ApiKey "+apikey)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 201 {
		t.Fatal(response.Status)
	}

	buf, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var feed database.Feed
	err = json.Unmarshal(buf, &feed)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFeeds(t *testing.T) {
	request, _ := http.NewRequest("GET", apiAddr+"/feeds", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal(response.Status)
	}

	buf, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var feeds []database.Feed
	err = json.Unmarshal(buf, &feeds)
	if err != nil {
		t.Fatal(err)
	}

	//println(string(buf))
}
