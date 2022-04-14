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

func TestCreateChildDownload(t *testing.T) {

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
		inputJSON     string
		statusCode    int
		targetPath    string
		receivedBytes uint64
		totalBytes    uint64
		siteUrl       string
		tabUrl        string
		mimeType      string
		child_id      uint64
		tokenGiven    string
		errorMessage  string
	}{
		{
			inputJSON:     `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 1}`,
			statusCode:    201,
			tokenGiven:    tokenString,
			child_id:      child.ID,
			targetPath:    "D:/",
			receivedBytes: 100,
			totalBytes:    100,
			siteUrl:       "www.google.com",
			tabUrl:        "google.com/tabURL",
			mimeType:      "text/html",
			errorMessage:  "",
		},
		{
			// When no token is passed
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"target_path": "", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh target_path",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh site_url",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "", "mime_type": "text/html", "child_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh tab_url",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "", "child_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh mime_type",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 0}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh child",
		},
		{
			// When child 2 uses child 1 token
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "child_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/childdownloads", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateChildDownload)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["child_id"], float64(v.child_id)) //just for both ids to have the same type
			assert.Equal(t, responseMap["target_path"], v.targetPath)
			assert.Equal(t, responseMap["received_bytes"], float64(v.receivedBytes))
			assert.Equal(t, responseMap["total_bytes"], float64(v.totalBytes))
			assert.Equal(t, responseMap["site_url"], v.siteUrl)
			assert.Equal(t, responseMap["tab_url"], v.tabUrl)
			assert.Equal(t, responseMap["mime_type"], v.mimeType)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetChildDownloads(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedParentsAndChilds()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedChildDownloads()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/childdownloads", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetChildDownloads)
	handler.ServeHTTP(rr, req)

	var childdownloads []models.ChildDownload
	_ = json.Unmarshal(rr.Body.Bytes(), &childdownloads)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(childdownloads), 2)
}
func TestGetChildDownloadByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatal(err)
	}
	childDownload, err := seedOneChildDownload()
	if err != nil {
		log.Fatal(err)
	}
	childDownload.Child = child
	childDownload.ChildID = child.ID
	childDownloadSample := []struct {
		id             string
		statusCode     int
		targetPath     string
		receivedBytes  int
		totalBytes     int
		siteUrl        string
		tabUrl         string
		tabReferredUrl string
		mimeType       string
		child_id       int
		errorMessage   string
	}{
		{
			id:             strconv.Itoa(int(childDownload.ID)),
			statusCode:     200,
			targetPath:     childDownload.TargetPath,
			receivedBytes:  int(childDownload.ReceivedBytes),
			totalBytes:     int(childDownload.TotalBytes),
			siteUrl:        childDownload.SiteUrl,
			tabUrl:         childDownload.TabUrl,
			tabReferredUrl: childDownload.TabReferredUrl,
			mimeType:       childDownload.MimeType,
			child_id:       int(childDownload.ChildID),
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range childDownloadSample {

		req, err := http.NewRequest("GET", "/childdownloads", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetChildDownload)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, float64(childDownload.ID), responseMap["id"])
			assert.Equal(t, childDownload.TargetPath, responseMap["target_path"]) //the response url id is float64
			assert.Equal(t, float64(childDownload.ReceivedBytes), responseMap["received_bytes"])
			assert.Equal(t, float64(childDownload.TotalBytes), responseMap["total_bytes"])
			assert.Equal(t, childDownload.SiteUrl, responseMap["site_url"])
			assert.Equal(t, childDownload.TabUrl, responseMap["tab_url"])
			assert.Equal(t, childDownload.TabReferredUrl, responseMap["tab_referred_url"])
			assert.Equal(t, childDownload.MimeType, responseMap["mime_type"])
			assert.Equal(t, float64(childDownload.ChildID), responseMap["child_id"]) //the response child id is float64
		}
	}
}

func TestDeleteChildDownload(t *testing.T) {
	var ParentEmail, ParentPassword string
	var AuthChilddownloadID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, _, childdownloads, err := seedParentsAndChildsAndChildDownloads()
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

	// Get only the second childdownload
	for _, childdownload := range childdownloads {
		if childdownload.ID == 1 {
			continue
		}
		AuthChilddownloadID = childdownload.ID
	}
	childdownloadSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthChilddownloadID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthChilddownloadID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthChilddownloadID)),
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
	for _, v := range childdownloadSample {

		req, _ := http.NewRequest("GET", "/childdownloads", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteChildDownload)

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
