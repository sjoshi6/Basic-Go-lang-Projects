package main

import (
	"fmt"
	"go-lbapp/api"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server          *httptest.Server
	reader          io.Reader
	searcheventsURL string
)

func init() {

	router := api.GetRouter()
	server = httptest.NewServer(router) //Creating new server with the user handlers
	searcheventsURL = fmt.Sprintf("%s/v1/search_events", server.URL)
}

// TestSearchEvents : Check if search events route works
func TestSearchEvents(t *testing.T) {
	searchJSON := `{"latitude":"100.2", "longitude":"127.2", "radius": "5"}`

	reader = strings.NewReader(searchJSON) //Convert string to reader

	request, err := http.NewRequest("GET", searcheventsURL, reader) //Create request with JSON body

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
	}
}
