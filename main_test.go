package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const (
	apiAddr = "http://localhost:8080/v1"
)

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
	} else if response.StatusCode != 200 {
		t.Fatal(response.Status)
	}

	result, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	println(string(result))
}
