package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
)

func TestCreateUrl(t *testing.T) {
	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		url          string
		title        string
		errorMessage string
	}{
		{
			inputJSON:    `{"url":"www.google.com", "title": "Google"}`,
			statusCode:   201,
			url:          "www.google.com",
			title:        "Google",
			errorMessage: "",
		},
		{
			inputJSON:    `{"url":"", "title": "Google"}`,
			statusCode:   422,
			errorMessage: "butuh url",
		},
		{
			inputJSON:    `{"url":"www.google.com", "title": ""}`,
			statusCode:   422,
			errorMessage: "butuh title",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/urls", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUrl)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["url"], v.url)
			assert.Equal(t, responseMap["title"], v.title)
		}
		if v.statusCode == 422 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUrls(t *testing.T) {
	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	err = seedUrls()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/urls", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUrls)
	handler.ServeHTTP(rr, req)

	var urls []models.Url
	err = json.Unmarshal(rr.Body.Bytes(), &urls)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(urls), 2)
}

func TestGetUrlByID(t *testing.T) {
	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatal(err)
	}
	urlSample := []struct {
		id           string
		statusCode   int
		url          string
		title        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(url.ID)),
			statusCode: 200,
			url:        url.Url,
			title:      url.Title,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range urlSample {

		req, err := http.NewRequest("GET", "/urls", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUrl)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, url.Url, responseMap["url"])
			assert.Equal(t, url.Title, responseMap["title"])
		}
	}
}
