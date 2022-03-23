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

func TestCreateParent(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		nama         string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"nama":"Pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			nama:         "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"nama":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "email sudah diambil",
		},
		{
			inputJSON:    `{"nama":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    `{"nama": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "butuh nama",
		},
		{
			inputJSON:    `{"nama": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "butuh email",
		},
		{
			inputJSON:    `{"nama": "Kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "butuh password",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parents", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateParent)
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
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetParents(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedParents()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/parents", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetParents)
	handler.ServeHTTP(rr, req)

	var parents []models.Parent
	err = json.Unmarshal(rr.Body.Bytes(), &parents)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(parents), 2)
}

func TestGetParentByID(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatal(err)
	}
	parentSample := []struct {
		id           string
		statusCode   int
		nama         string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(parent.ID)),
			statusCode: 200,
			nama:       parent.Nama,
			email:      parent.Email,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range parentSample {

		req, err := http.NewRequest("GET", "/parents", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetParent)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, parent.Nama, responseMap["nama"])
			assert.Equal(t, parent.Email, responseMap["email"])
		}
	}
}

func TestUpdateParent(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, err := seedParents() //we need atleast two parents to properly check the update
	if err != nil {
		log.Fatalf("Error seeding parent: %v\n", err)
	}
	// Get only the first parent
	for _, parent := range parents {
		if parent.ID == 2 {
			continue
		}
		AuthID = parent.ID
		AuthEmail = parent.Email
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the parent and get the authentication token
	token, err := server.ParentSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login parent: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		updateNama   string
		updateEmail  string
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Grand", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   200,
			updateNama:   "Grand",
			updateEmail:  "grand@gmail.com",
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Woman", "email": "woman@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Man", "email": "man@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Woman", "email": "woman@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to parent 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Frank", "email": "kenny@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email sudah diambil",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "invalid email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh nama",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh email",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When parent 2 is using parent 1 token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"nama": "Mike", "email": "mike@gmail.com", "password": "password"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parents", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateParent)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nama"], v.updateNama)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateParentProfile(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, err := seedParents() //we need atleast two parents to properly check the update
	if err != nil {
		log.Fatalf("Error seeding parent: %v\n", err)
	}
	// Get only the first parent
	for _, parent := range parents {
		if parent.ID == 2 {
			continue
		}
		AuthID = parent.ID
		AuthEmail = parent.Email
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the parent and get the authentication token
	token, err := server.ParentSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login parent: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		updateNama   string
		updateEmail  string
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Grand", "email": "grand@gmail.com"}`,
			statusCode:   200,
			updateNama:   "Grand",
			updateEmail:  "grand@gmail.com",
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Man", "email": "man@gmail.com"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Woman", "email": "woman@gmail.com"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to parent 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Frank", "email": "kenny@gmail.com"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email sudah diambil",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama":"Kan", "email": "kangmail.com"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "invalid email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama": "", "email": "kan@gmail.com"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh nama",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nama": "Kan", "email": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh email",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When parent 2 is using parent 1 token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"nama": "Mike", "email": "mike@gmail.com"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parents", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateParentProfile)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nama"], v.updateNama)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateParentPassword(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, err := seedParents() //we need atleast two parents to properly check the update
	if err != nil {
		log.Fatalf("Error seeding parent: %v\n", err)
	}
	// Get only the first parent
	for _, parent := range parents {
		if parent.ID == 2 {
			continue
		}
		AuthID = parent.ID
		AuthEmail = parent.Email
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the parent and get the authentication token
	token, err := server.ParentSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login parent: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updatePassword string
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to string
			id:             strconv.Itoa(int(AuthID)),
			updateJSON:     `{"oldPassword": "password", "newPassword": "newpassword"}`,
			statusCode:     200,
			updatePassword: "newpassword",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"oldPassword": "", "newPassword": "newpassword"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh old password",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"oldPassword": "password", "newPassword": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh new password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"oldPassword": "password", "newPassword": "newpassword"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"oldPassword": "password", "newPassword": "newpassword"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parents/password", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateParentPassword)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			parent := models.Parent{}
			err = server.DB.Debug().Model(models.Parent{}).Where("id = ?", v.id).Take(&parent).Error
			if err != nil {
				assert.Equal(t, err, errors.New(v.errorMessage))
			}
			token, err := server.ParentSignIn(parent.Email, v.updatePassword)
			if err != nil {
				assert.Equal(t, err, errors.New(v.errorMessage))
			} else {
				assert.NotEqual(t, token, "")
			}
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteParent(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}

	parents, err := seedParents() //we need atleast two parents to properly check the update
	if err != nil {
		log.Fatalf("Error seeding parent: %v\n", err)
	}
	// Get only the first and log him in
	for _, parent := range parents {
		if parent.ID == 2 {
			continue
		}
		AuthID = parent.ID
		AuthEmail = parent.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}
	//Login the parent and get the authentication token
	token, err := server.ParentSignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	parentSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
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
			// Parent 2 trying to use Parent 1 token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range parentSample {

		req, err := http.NewRequest("GET", "/parents", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteParent)

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
