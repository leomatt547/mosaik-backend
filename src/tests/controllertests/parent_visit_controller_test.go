package controllertests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
)

// func TestCreateParent(t *testing.T) {

// 	err := refreshParentAndParentTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	parent, err := seedOneParent()
// 	if err != nil {
// 		log.Fatalf("Cannot seed parent %v\n", err)
// 	}
// 	token, err := server.SignIn(parent.Email, "password") //Note the password in the database is already hashed, we want unhashed
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	samples := []struct {
// 		inputJSON    string
// 		statusCode   int
// 		nama        string
// 		email      string
// 		parent_id    uint32
// 		tokenGiven   string
// 		errorMessage string
// 	}{
// 		{
// 			inputJSON:    `{"nama":"The nama", "email": "the email", "parent_id": 1}`,
// 			statusCode:   201,
// 			tokenGiven:   tokenString,
// 			nama:        "The nama",
// 			email:      "the email",
// 			parent_id:    parent.ID,
// 			errorMessage: "",
// 		},
// 		{
// 			inputJSON:    `{"nama":"The nama", "email": "the email", "parent_id": 1}`,
// 			statusCode:   500,
// 			tokenGiven:   tokenString,
// 			errorMessage: "nama Already Taken",
// 		},
// 		{
// 			// When no token is passed
// 			inputJSON:    `{"nama":"When no token is passed", "email": "the email", "parent_id": 1}`,
// 			statusCode:   401,
// 			tokenGiven:   "",
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			// When incorrect token is passed
// 			inputJSON:    `{"nama":"When incorrect token is passed", "email": "the email", "parent_id": 1}`,
// 			statusCode:   401,
// 			tokenGiven:   "This is an incorrect token",
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			inputJSON:    `{"nama": "", "email": "The email", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "butuh nama",
// 		},
// 		{
// 			inputJSON:    `{"nama": "This is a nama", "email": "", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "butuh email",
// 		},
// 		{
// 			inputJSON:    `{"nama": "This is an awesome nama", "email": "the email"}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "butuh Author",
// 		},
// 		{
// 			// When parent 2 uses parent 1 token
// 			inputJSON:    `{"nama": "This is an awesome nama", "email": "the email", "parent_id": 2}`,
// 			statusCode:   401,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Unauthorized",
// 		},
// 	}
// 	for _, v := range samples {

// 		req, err := http.NewRequest("POST", "/parents", bytes.NewBufferString(v.inputJSON))
// 		if err != nil {
// 			t.Errorf("this is the error: %v\n", err)
// 		}
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.CreateParent)

// 		req.Header.Set("Authorization", v.tokenGiven)
// 		handler.ServeHTTP(rr, req)

// 		responseMap := make(map[string]interface{})
// 		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
// 		if err != nil {
// 			fmt.Printf("Cannot convert to json: %v", err)
// 		}
// 		assert.Equal(t, rr.Code, v.statusCode)
// 		if v.statusCode == 201 {
// 			assert.Equal(t, responseMap["nama"], v.nama)
// 			assert.Equal(t, responseMap["email"], v.email)
// 			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just for both ids to have the same type
// 		}
// 		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
// 			assert.Equal(t, responseMap["error"], v.errorMessage)
// 		}
// 	}
// }

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

// func TestUpdateParent(t *testing.T) {

// 	var ParentParentEmail, ParentParentPassword string
// 	var AuthParentParentID uint32
// 	var AuthParentID uint64

// 	err := refreshParentAndParentTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	parents, parentvisits, err := seedParents()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Get only the first parent
// 	for _, parent := range parents {
// 		if parent.ID == 2 {
// 			continue
// 		}
// 		ParentParentEmail = parent.Email
// 		ParentParentPassword = "password" //Note the password in the database is already hashed, we want unhashed
// 	}
// 	//Login the parent and get the authentication token
// 	token, err := server.SignIn(ParentParentEmail, ParentParentPassword)
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	// Get only the first parent
// 	for _, parent := range parents {
// 		if parent.ID == 2 {
// 			continue
// 		}
// 		AuthParentID = parent.ID
// 		AuthParentParentID = parent.ParentID
// 	}
// 	// fmt.Printf("this is the auth parent: %v\n", AuthParentID)

// 	samples := []struct {
// 		id           string
// 		updateJSON   string
// 		statusCode   int
// 		nama        string
// 		email      string
// 		parent_id    uint32
// 		tokenGiven   string
// 		errorMessage string
// 	}{
// 		{
// 			// Convert int64 to int first before converting to string
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"The updated parent", "email": "This is the updated email", "parent_id": 1}`,
// 			statusCode:   200,
// 			nama:        "The updated parent",
// 			email:      "This is the updated email",
// 			parent_id:    AuthParentParentID,
// 			tokenGiven:   tokenString,
// 			errorMessage: "",
// 		},
// 		{
// 			// When no token is provided
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"This is still another nama", "email": "This is the updated email", "parent_id": 1}`,
// 			tokenGiven:   "",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			// When incorrect token is provided
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"This is still another nama", "email": "This is the updated email", "parent_id": 1}`,
// 			tokenGiven:   "this is an incorrect token",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			//Note: "nama 2" belongs to parent 2, and nama must be unique
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"nama 2", "email": "This is the updated email", "parent_id": 1}`,
// 			statusCode:   500,
// 			tokenGiven:   tokenString,
// 			errorMessage: "nama Already Taken",
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"", "email": "This is the updated email", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "butuh nama",
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"Awesome nama", "email": "", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "butuh email",
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"This is another nama", "email": "This is the updated email"}`,
// 			statusCode:   401,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			id:         "unknown",
// 			statusCode: 400,
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			updateJSON:   `{"nama":"This is still another nama", "email": "This is the updated email", "parent_id": 2}`,
// 			tokenGiven:   tokenString,
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 	}

// 	for _, v := range samples {

// 		req, err := http.NewRequest("POST", "/parents", bytes.NewBufferString(v.updateJSON))
// 		if err != nil {
// 			t.Errorf("this is the error: %v\n", err)
// 		}
// 		req = mux.SetURLVars(req, map[string]string{"id": v.id})
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.UpdateParent)

// 		req.Header.Set("Authorization", v.tokenGiven)

// 		handler.ServeHTTP(rr, req)

// 		responseMap := make(map[string]interface{})
// 		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
// 		if err != nil {
// 			t.Errorf("Cannot convert to json: %v", err)
// 		}
// 		assert.Equal(t, rr.Code, v.statusCode)
// 		if v.statusCode == 200 {
// 			assert.Equal(t, responseMap["nama"], v.nama)
// 			assert.Equal(t, responseMap["email"], v.email)
// 			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just to match the type of the json we receive thats why we used float64
// 		}
// 		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
// 			assert.Equal(t, responseMap["error"], v.errorMessage)
// 		}
// 	}
// }

// func TestDeleteParent(t *testing.T) {

// 	var ParentParentEmail, ParentParentPassword string
// 	var ParentParentID uint32
// 	var AuthParentID uint64

// 	err := refreshParentAndParentTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	parents, parents, err := seedParents()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//Let's get only the Second parent
// 	for _, parent := range parents {
// 		if parent.ID == 1 {
// 			continue
// 		}
// 		ParentParentEmail = parent.Email
// 		ParentParentPassword = "password" //Note the password in the database is already hashed, we want unhashed
// 	}
// 	//Login the parent and get the authentication token
// 	token, err := server.SignIn(ParentParentEmail, ParentParentPassword)
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	// Get only the second parent
// 	for _, parent := range parents {
// 		if parent.ID == 1 {
// 			continue
// 		}
// 		AuthParentID = parent.ID
// 		ParentParentID = parent.ParentID
// 	}
// 	parentSample := []struct {
// 		id           string
// 		parent_id    uint32
// 		tokenGiven   string
// 		statusCode   int
// 		errorMessage string
// 	}{
// 		{
// 			// Convert int64 to int first before converting to string
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			parent_id:    ParentParentID,
// 			tokenGiven:   tokenString,
// 			statusCode:   204,
// 			errorMessage: "",
// 		},
// 		{
// 			// When empty token is passed
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			parent_id:    ParentParentID,
// 			tokenGiven:   "",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			// When incorrect token is passed
// 			id:           strconv.Itoa(int(AuthParentID)),
// 			parent_id:    ParentParentID,
// 			tokenGiven:   "This is an incorrect token",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			id:         "unknwon",
// 			tokenGiven: tokenString,
// 			statusCode: 400,
// 		},
// 		{
// 			id:           strconv.Itoa(int(1)),
// 			parent_id:    1,
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 	}
// 	for _, v := range parentSample {

// 		req, _ := http.NewRequest("GET", "/parents", nil)
// 		req = mux.SetURLVars(req, map[string]string{"id": v.id})

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.DeleteParent)

// 		req.Header.Set("Authorization", v.tokenGiven)

// 		handler.ServeHTTP(rr, req)

// 		assert.Equal(t, rr.Code, v.statusCode)

// 		if v.statusCode == 401 && v.errorMessage != "" {

// 			responseMap := make(map[string]interface{})
// 			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
// 			if err != nil {
// 				t.Errorf("Cannot convert to json: %v", err)
// 			}
// 			assert.Equal(t, responseMap["error"], v.errorMessage)
// 		}
// 	}
// }
