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

func TestCreateNSFWUrl(t *testing.T) {
	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		url          string
		errorMessage string
	}{
		{
			inputJSON:    `{"url":"www.pornhub.com"}`,
			statusCode:   201,
			url:          "www.pornhub.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"url":""}`,
			statusCode:   422,
			errorMessage: "butuh url",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/nsfw", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateNSFWUrl)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["url"], v.url)
		}
		if v.statusCode == 422 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetNSFWUrls(t *testing.T) {
	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	err = seedNSFWUrls()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/nsfw", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetNSFWUrls)
	handler.ServeHTTP(rr, req)

	var nsfw []models.NSFWUrl
	err = json.Unmarshal(rr.Body.Bytes(), &nsfw)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(nsfw), 2)
}

func TestGetNSFWUrlByID(t *testing.T) {
	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	nsfw, err := seedOneNSFWUrl()
	if err != nil {
		log.Fatal(err)
	}
	nsfwSample := []struct {
		id           string
		statusCode   int
		url          string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(nsfw.ID)),
			statusCode: 200,
			url:        nsfw.Url,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range nsfwSample {

		req, err := http.NewRequest("GET", "/nsfw", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetNSFWUrl)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, nsfw.Url, responseMap["url"])
		}
	}
}

// func TestSavedSearchChecker(t *testing.T) {
// 	err := refreshNSFWUrlTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	samples := []struct {
// 		inputJSON    string
// 		statusCode   int
// 		url          string
// 		isBlocked    bool
// 		errorMessage string
// 	}{
// 		{
// 			inputJSON:    `{"url":"https://m.jpnn.com/news/4-potret-seksi-dea-onlyfans-pakai-lingerie-hingga-bralette"}`,
// 			statusCode:   200,
// 			url:          "https://m.jpnn.com/news/4-potret-seksi-dea-onlyfans-pakai-lingerie-hingga-bralette",
// 			isBlocked:    true,
// 			errorMessage: "",
// 		},
// 	}

// 	for _, v := range samples {
// 		req, err := http.NewRequest("POST", "/nsfw", bytes.NewBufferString(v.inputJSON))
// 		if err != nil {
// 			t.Errorf("this is the error: %v", err)
// 		}
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.SavedSearchChecker)
// 		handler.ServeHTTP(rr, req)

// 		responseMap := make(map[string]interface{})
// 		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
// 		if err != nil {
// 			fmt.Printf("Cannot convert to json: %v", err)
// 		}
// 		assert.Equal(t, rr.Code, v.statusCode)
// 		if v.statusCode == 200 {
// 			assert.Equal(t, responseMap["url"], v.url)
// 			assert.Equal(t, responseMap["is_blocked"], v.isBlocked)
// 		}
// 	}
// }
