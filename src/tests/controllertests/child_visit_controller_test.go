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

// func TestCreateChild(t *testing.T) {

// 	err := refreshParentAndChildTable()
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
// 			errorMessage: "Nama Already Taken",
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
// 			errorMessage: "Required Nama",
// 		},
// 		{
// 			inputJSON:    `{"nama": "This is a nama", "email": "", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Required Email",
// 		},
// 		{
// 			inputJSON:    `{"nama": "This is an awesome nama", "email": "the email"}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Required Author",
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

// 		req, err := http.NewRequest("POST", "/childs", bytes.NewBufferString(v.inputJSON))
// 		if err != nil {
// 			t.Errorf("this is the error: %v\n", err)
// 		}
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.CreateChild)

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

// func TestUpdateChild(t *testing.T) {

// 	var ChildParentEmail, ChildParentPassword string
// 	var AuthChildParentID uint32
// 	var AuthChildID uint64

// 	err := refreshParentAndChildTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	parents, childvisits, err := seedParentsAndChilds()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Get only the first parent
// 	for _, parent := range parents {
// 		if parent.ID == 2 {
// 			continue
// 		}
// 		ChildParentEmail = parent.Email
// 		ChildParentPassword = "password" //Note the password in the database is already hashed, we want unhashed
// 	}
// 	//Login the parent and get the authentication token
// 	token, err := server.SignIn(ChildParentEmail, ChildParentPassword)
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	// Get only the first child
// 	for _, child := range childs {
// 		if child.ID == 2 {
// 			continue
// 		}
// 		AuthChildID = child.ID
// 		AuthChildParentID = child.ParentID
// 	}
// 	// fmt.Printf("this is the auth child: %v\n", AuthChildID)

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
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"The updated child", "email": "This is the updated email", "parent_id": 1}`,
// 			statusCode:   200,
// 			nama:        "The updated child",
// 			email:      "This is the updated email",
// 			parent_id:    AuthChildParentID,
// 			tokenGiven:   tokenString,
// 			errorMessage: "",
// 		},
// 		{
// 			// When no token is provided
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"This is still another nama", "email": "This is the updated email", "parent_id": 1}`,
// 			tokenGiven:   "",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			// When incorrect token is provided
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"This is still another nama", "email": "This is the updated email", "parent_id": 1}`,
// 			tokenGiven:   "this is an incorrect token",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			//Note: "Nama 2" belongs to child 2, and nama must be unique
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"Nama 2", "email": "This is the updated email", "parent_id": 1}`,
// 			statusCode:   500,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Nama Already Taken",
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"", "email": "This is the updated email", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Required Nama",
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"Awesome nama", "email": "", "parent_id": 1}`,
// 			statusCode:   422,
// 			tokenGiven:   tokenString,
// 			errorMessage: "Required Email",
// 		},
// 		{
// 			id:           strconv.Itoa(int(AuthChildID)),
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
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			updateJSON:   `{"nama":"This is still another nama", "email": "This is the updated email", "parent_id": 2}`,
// 			tokenGiven:   tokenString,
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 	}

// 	for _, v := range samples {

// 		req, err := http.NewRequest("POST", "/childs", bytes.NewBufferString(v.updateJSON))
// 		if err != nil {
// 			t.Errorf("this is the error: %v\n", err)
// 		}
// 		req = mux.SetURLVars(req, map[string]string{"id": v.id})
// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.UpdateChild)

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

// func TestDeleteChild(t *testing.T) {

// 	var ChildParentEmail, ChildParentPassword string
// 	var ChildParentID uint32
// 	var AuthChildID uint64

// 	err := refreshParentAndChildTable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	parents, childs, err := seedParentsAndChilds()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//Let's get only the Second parent
// 	for _, parent := range parents {
// 		if parent.ID == 1 {
// 			continue
// 		}
// 		ChildParentEmail = parent.Email
// 		ChildParentPassword = "password" //Note the password in the database is already hashed, we want unhashed
// 	}
// 	//Login the parent and get the authentication token
// 	token, err := server.SignIn(ChildParentEmail, ChildParentPassword)
// 	if err != nil {
// 		log.Fatalf("cannot login: %v\n", err)
// 	}
// 	tokenString := fmt.Sprintf("Bearer %v", token)

// 	// Get only the second child
// 	for _, child := range childs {
// 		if child.ID == 1 {
// 			continue
// 		}
// 		AuthChildID = child.ID
// 		ChildParentID = child.ParentID
// 	}
// 	childSample := []struct {
// 		id           string
// 		parent_id    uint32
// 		tokenGiven   string
// 		statusCode   int
// 		errorMessage string
// 	}{
// 		{
// 			// Convert int64 to int first before converting to string
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			parent_id:    ChildParentID,
// 			tokenGiven:   tokenString,
// 			statusCode:   204,
// 			errorMessage: "",
// 		},
// 		{
// 			// When empty token is passed
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			parent_id:    ChildParentID,
// 			tokenGiven:   "",
// 			statusCode:   401,
// 			errorMessage: "Unauthorized",
// 		},
// 		{
// 			// When incorrect token is passed
// 			id:           strconv.Itoa(int(AuthChildID)),
// 			parent_id:    ChildParentID,
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
// 	for _, v := range childSample {

// 		req, _ := http.NewRequest("GET", "/childs", nil)
// 		req = mux.SetURLVars(req, map[string]string{"id": v.id})

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(server.DeleteChild)

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