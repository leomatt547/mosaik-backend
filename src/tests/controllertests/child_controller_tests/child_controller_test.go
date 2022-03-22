package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestCreateChild(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Cannot seed parent %v\n", err)
	}
	token, err := server.ParentSignIn(parent.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		nama         string
		email        string
		parent_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"nama":"jr", "email": "jr@gmail.com", "password":"jr123" , "parent_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			nama:         "jr",
			email:        "jr@gmail.com",
			parent_id:    parent.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"nama":"jr_2", "email": "jr@gmail.com",  "password":"jr123", "parent_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email sudah diambil",
		},
		{
			// When no token is passed
			inputJSON:    `{"nama":"When no token is passed", "email": "jr@gmail.com",  "password":"jr123" ,"parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"nama":"When incorrect token is passed", "email": "jr@gmail.com", "password":"jr123" , "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"nama": "", "email": "jr@gmail.com", "password":"jr123" , "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh nama",
		},
		{
			inputJSON:    `{"nama": "This is a nama", "email": "", "password":"jr123" , "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh email",
		},
		{
			inputJSON:    `{"nama": "This is an awesome nama", "email": "jr@gmail.com", "password":"jr123"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh parent_id",
		},
		{
			inputJSON:    `{"nama": "This is an awesome nama", "email": "jr@gmail.com", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh password",
		},
		{
			// When parent 2 uses parent 1 token
			inputJSON:    `{"nama": "This is an awesome nama", "email": "jr@gmail.com", "password":"jr123", "parent_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/childs", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateChild)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["nama"], v.nama)
			assert.Equal(t, responseMap["email"], v.email)
			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetChilds(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedParentsAndChilds()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/childs", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetChilds)
	handler.ServeHTTP(rr, req)

	var childs []models.Child
	_ = json.Unmarshal(rr.Body.Bytes(), &childs)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(childs), 2)
}
func TestGetChildByID(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatal(err)
	}
	childSample := []struct {
		id           string
		statusCode   int
		nama         string
		email        string
		parent_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(child.ID)),
			statusCode: 200,
			nama:       child.Nama,
			email:      child.Email,
			parent_id:  child.ParentID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range childSample {

		req, err := http.NewRequest("GET", "/childs", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetChild)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, child.Nama, responseMap["nama"])
			assert.Equal(t, child.Email, responseMap["email"])
			assert.Equal(t, float64(child.ParentID), responseMap["parent_id"]) //the response author id is float64
		}
	}
}

func TestUpdateChild(t *testing.T) {

	var ParentEmail, ParentPassword string
	var AuthParentID uint32
	var AuthChildID uint64

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, childs, err := seedParentsAndChilds()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first parent
	for _, parent := range parents {
		if parent.ID == 2 {
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

	// Get only the first child
	for _, child := range childs {
		if child.ID == 2 {
			continue
		}
		AuthChildID = child.ID
		AuthParentID = child.ParentID
	}
	// fmt.Printf("this is the auth child: %v\n", AuthChildID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		nama         string
		email        string
		parent_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"jr_2", "email": "jr_2@gmail.com", "password":"jr123", "parent_id": 1}`,
			statusCode:   200,
			nama:         "jr_2",
			email:        "jr_2@gmail.com",
			parent_id:    AuthParentID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"This is still another nama", "email": "jr_2@gmail.com", "password":"jr123", "parent_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"This is still another nama", "email": "jr_2@gmail.com", "password":"jr123", "parent_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Magu Frank" belongs to child 2, and nama must be unique
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"Martin Luth Junior", "email": "magu_jr@gmail.com", "password":"jr123", "parent_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email sudah diambil",
		},
		{
			id:           strconv.Itoa(int(AuthChildID + 1)),
			updateJSON:   `{"nama":"This is another nama", "password":"jr123", "email": "jr_2@gmail.com"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthChildID + 1)),
			updateJSON:   `{"nama":"This is still another nama", "email": "jr_2@gmail.com", "password":"jr123", "parent_id": 1}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/childs", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateChild)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nama"], v.nama)
			assert.Equal(t, responseMap["email"], v.email)
			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateChildProfile(t *testing.T) {

	var ParentEmail, ParentPassword string
	var AuthParentID uint32
	var AuthChildID uint64

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, childs, err := seedParentsAndChilds()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first parent
	for _, parent := range parents {
		if parent.ID == 2 {
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

	// Get only the first child
	for _, child := range childs {
		if child.ID == 2 {
			continue
		}
		AuthChildID = child.ID
		AuthParentID = child.ParentID
	}
	// fmt.Printf("this is the auth child: %v\n", AuthChildID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		nama         string
		email        string
		parent_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"jr_2", "email": "jr_2@gmail.com", "parent_id": 1}`,
			statusCode:   200,
			nama:         "jr_2",
			email:        "jr_2@gmail.com",
			parent_id:    AuthParentID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"This is still another nama", "email": "jr_2@gmail.com", "parent_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"This is still another nama", "email": "jr_2@gmail.com", "parent_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Magu Frank" belongs to child 2, and nama must be unique
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"nama":"Martin Luth Junior", "email": "magu_jr@gmail.com", "parent_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email sudah diambil",
		},
		{
			id:           strconv.Itoa(int(AuthChildID + 1)),
			updateJSON:   `{"nama":"This is another nama", "email": "jr_2@gmail.com"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthChildID + 1)),
			updateJSON:   `{"nama":"This is still another nama", "email": "jr_2@gmail.com", "parent_id": 1}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/childs", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateChildProfile)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nama"], v.nama)
			assert.Equal(t, responseMap["email"], v.email)
			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateChildPassword(t *testing.T) {

	var ParentEmail, ParentPassword string
	var AuthParentID uint32
	var AuthChildID uint64

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, childs, err := seedParentsAndChilds() //we need atleast two parents to properly check the update
	if err != nil {
		log.Fatalf("Error seeding parent: %v\n", err)
	}
	// Get only the first parent
	for _, parent := range parents {
		if parent.ID == 2 {
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

	// Get only the first child
	for _, child := range childs {
		if child.ID == 2 {
			continue
		}
		AuthChildID = child.ID
		AuthParentID = child.ParentID
	}

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		email          string
		updatePassword string
		parent_id      uint32
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to string
			id:             strconv.Itoa(int(AuthChildID)),
			updateJSON:     `{"email": "steven_jr@gmail.com", "oldPassword": "password", "newPassword": "newpassword", "parent_id": 1}`,
			statusCode:     200,
			email:          "steven_jr@gmail.com",
			updatePassword: "newpassword",
			parent_id:      AuthParentID,
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"email": "steven_jr@gmail.com", "oldPassword": "", "newPassword": "newpassword", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh old password",
		},
		{
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"email": "steven_jr@gmail.com", "oldPassword": "password", "newPassword": "", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh new password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"email": "steven_jr@gmail.com", "oldPassword": "password", "newPassword": "newpassword", "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"email": "steven_jr@gmail.com", "oldPassword": "password", "newPassword": "newpassword",  "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(AuthChildID)),
			updateJSON:   `{"email": "", "oldPassword": "password", "newPassword": "newpassword", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh email",
		},
		{
			id:           strconv.Itoa(int(AuthChildID + 1)),
			updateJSON:   `{"email": "steven_jr@gmail.com", "oldPassword": "password", "newPassword": "newpassword"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthChildID + 1)),
			updateJSON:   `{"email": "steven_jr@gmail.com", "oldPassword": "password", "newPassword": "newpassword",  "parent_id": 1}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/childs/password", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateChildPassword)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			token, err := server.ChildSignIn(v.email, v.updatePassword)
			if err != nil {
				assert.Equal(t, err, errors.New(v.errorMessage))
			} else {
				assert.NotEqual(t, token, "")
				assert.Equal(t, responseMap["parent_id"], float64(v.parent_id))
			}
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteChild(t *testing.T) {

	var ParentEmail, ParentPassword string
	var ParentID uint32
	var AuthChildID uint64

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, childs, err := seedParentsAndChilds()
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

	// Get only the second child
	for _, child := range childs {
		if child.ID == 1 {
			continue
		}
		AuthChildID = child.ID
		ParentID = child.ParentID
	}
	childSample := []struct {
		id           string
		parent_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthChildID)),
			parent_id:    ParentID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthChildID)),
			parent_id:    ParentID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthChildID)),
			parent_id:    ParentID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
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
	for _, v := range childSample {

		req, _ := http.NewRequest("GET", "/childs", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteChild)

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
