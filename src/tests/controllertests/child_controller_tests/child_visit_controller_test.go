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

func TestCreateChildVisit(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	child, _, err := seedOneParentAndOneChildAndOneUrl()
	if err != nil {
		log.Fatalf("Cannot seed child %v\n", err)
	}
	response, err := server.ChildSignIn(child.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", response.Token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		url_id       uint64
		duration     uint64
		child_id     uint64
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"url_id": 1, "duration": 10, "child_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			url_id:       1,
			duration:     10,
			child_id:     child.ID,
			errorMessage: "",
		},
		{
			// When no token is passed
			inputJSON:    `{"url_id": 1, "duration": 10, "child_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"url_id": 1, "duration": 10, "child_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"url_id": 0, "duration": 10, "child_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh url_id",
		},
		{
			inputJSON:    `{"url_id": 1, "duration": 0, "child_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh durasi",
		},
		{
			inputJSON:    `{"url_id": 1, "duration": 10, "child_id": 0}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh child_id",
		},
		{
			// When child 2 uses child 1 token
			inputJSON:    `{"url_id": 1, "duration": 10, "child_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/childvisits", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateChildVisit)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["url_id"], float64(v.url_id))
			assert.Equal(t, responseMap["duration"], float64(v.duration))
			assert.Equal(t, responseMap["child_id"], float64(v.child_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetChildVisits(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedParentsAndChilds()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedChildVisitsAndUrls()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/childvisits", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetChildVisits)
	handler.ServeHTTP(rr, req)

	var childvisits []models.ChildVisit
	_ = json.Unmarshal(rr.Body.Bytes(), &childvisits)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(childvisits), 2)
}
func TestGetChildVisitByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatal(err)
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatal(err)
	}
	childVisit, err := seedOneChildVisit()
	if err != nil {
		log.Fatal(err)
	}
	childVisit.Child = child
	childVisit.ChildID = child.ID
	childVisit.Url = url
	childVisit.UrlID = url.ID
	childVisitSample := []struct {
		id           string
		statusCode   int
		url_id       int
		duration     int
		child_id     int
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(childVisit.ID)),
			statusCode: 200,
			url_id:     int(childVisit.UrlID),
			duration:   int(childVisit.Duration),
			child_id:   int(childVisit.ChildID),
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range childVisitSample {

		req, err := http.NewRequest("GET", "/childvisits", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetChildVisit)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, float64(childVisit.ID), responseMap["id"])
			assert.Equal(t, float64(childVisit.UrlID), responseMap["url_id"]) //the response url id is float64
			assert.Equal(t, float64(childVisit.Duration), responseMap["duration"])
			assert.Equal(t, float64(childVisit.ChildID), responseMap["child_id"]) //the response child id is float64
		}
	}
}

func TestDeleteChildVisit(t *testing.T) {
	var ParentEmail, ParentPassword string
	var AuthChildvisitID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, _, childvisits, _, err := seedParentsAndChildsAndChildVisitsAndUrls()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second parent
	for _, parent := range parents {
		if parent.ID == 1 {
			continue
		}
		ParentEmail = parent.Email
		ParentPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the parent and get the authentication token
	response, err := server.ParentSignIn(ParentEmail, ParentPassword, "")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", response.Token)

	// Get only the second childvisit
	for _, childvisit := range childvisits {
		if childvisit.ID == 1 {
			continue
		}
		AuthChildvisitID = childvisit.ID
	}
	childvisitSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthChildvisitID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthChildvisitID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthChildvisitID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range childvisitSample {

		req, _ := http.NewRequest("GET", "/childvisits", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteChildVisit)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
