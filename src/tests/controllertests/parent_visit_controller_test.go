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

func TestCreateParentVisit(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, _, err := seedParentsAndUrls()
	parent := parents[0]
	if err != nil {
		log.Fatalf("Cannot seed Parent and Urls %v\n", err)
	}
	token, err := server.ParentSignIn(parent.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		url_id       uint64
		duration     uint64
		parent_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"url_id": 1, "duration": 10, "parent_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			url_id:       1,
			duration:     10,
			parent_id:    parent.ID,
			errorMessage: "",
		},
		{
			// When no token is passed
			inputJSON:    `{"url_id": 1, "duration": 10, "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"url_id": 1, "duration": 10, "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"url_id": 0, "duration": 10, "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh url_id",
		},
		{
			inputJSON:    `{"url_id": 1, "duration": 0, "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh durasi",
		},
		{
			inputJSON:    `{"url_id": 1, "duration": 10, "parent_id": 0}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh parent_id",
		},
		{
			// When parent 2 uses parent 1 token
			inputJSON:    `{"url_id": 1, "duration": 10, "parent_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parentvisits", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateParentVisit)

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
			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetParentVisits(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedParents()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedParentVisitsAndUrls()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/parentvisits", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetParentVisits)
	handler.ServeHTTP(rr, req)

	var parentvisits []models.ParentVisit
	_ = json.Unmarshal(rr.Body.Bytes(), &parentvisits)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(parentvisits), 2)
}
func TestGetParentVisitByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatal(err)
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatal(err)
	}
	parentVisit, err := seedOneParentVisit()
	if err != nil {
		log.Fatal(err)
	}
	parentVisit.Parent = parent
	parentVisit.ParentID = parent.ID
	parentVisit.Url = url
	parentVisit.UrlID = url.ID
	parentVisitSample := []struct {
		id           string
		statusCode   int
		url_id       int
		duration     int
		parent_id    int
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(parentVisit.ID)),
			statusCode: 200,
			url_id:     int(parentVisit.UrlID),
			duration:   int(parentVisit.Duration),
			parent_id:  int(parentVisit.ParentID),
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range parentVisitSample {

		req, err := http.NewRequest("GET", "/parentvisits", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetParentVisit)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, float64(parentVisit.ID), responseMap["id"])
			assert.Equal(t, float64(parentVisit.UrlID), responseMap["url_id"]) //the response url id is float64
			assert.Equal(t, float64(parentVisit.Duration), responseMap["duration"])
			assert.Equal(t, float64(parentVisit.ParentID), responseMap["parent_id"]) //the response parent id is float64
		}
	}
}

func TestDeleteParentVisit(t *testing.T) {
	var ParentEmail, ParentPassword string
	var ParentID uint32
	var AuthParentvisitID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, parentvisits, _, err := seedParentsAndParentvisitsAndUrls()
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
	token, err := server.ParentSignIn(ParentEmail, ParentPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second parentvisit
	for _, parentvisit := range parentvisits {
		if parentvisit.ID == 1 {
			continue
		}
		AuthParentvisitID = parentvisit.ID
		ParentID = parentvisit.ParentID
	}
	parentvisitSample := []struct {
		id           string
		parent_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthParentvisitID)),
			parent_id:    ParentID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthParentvisitID)),
			parent_id:    ParentID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthParentvisitID)),
			parent_id:    ParentID,
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
			parent_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range parentvisitSample {

		req, _ := http.NewRequest("GET", "/parentvisits", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteParentVisit)

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
