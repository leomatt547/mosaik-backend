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

func TestCreateParentDownload(t *testing.T) {

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
		inputJSON     string
		statusCode    int
		targetPath    string
		receivedBytes uint64
		totalBytes    uint64
		siteUrl       string
		tabUrl        string
		mimeType      string
		parent_id     uint32
		tokenGiven    string
		errorMessage  string
	}{
		{
			inputJSON:     `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 1}`,
			statusCode:    201,
			tokenGiven:    tokenString,
			parent_id:     parent.ID,
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
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"target_path": "", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh target_path",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh site_url",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "", "mime_type": "text/html", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh tab_url",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "", "parent_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh mime_type",
		},
		{
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 0}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "butuh parent",
		},
		{
			// When parent 2 uses parent 1 token
			inputJSON:    `{"target_path": "D:/", "received_bytes": 100, "total_bytes": 100, "site_url":  "www.google.com", "tab_url": "google.com/tabURL", "mime_type": "text/html", "parent_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parentdownloads", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateParentDownload)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["parent_id"], float64(v.parent_id)) //just for both ids to have the same type
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

func TestGetParentDownloads(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedParentsAndDownloads()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/parentdownloads", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetParentDownloads)
	handler.ServeHTTP(rr, req)

	var parentdownloads []models.ParentDownload
	_ = json.Unmarshal(rr.Body.Bytes(), &parentdownloads)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(parentdownloads), 2)
}
func TestGetParentDownloadByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatal(err)
	}
	parentDownload, err := seedOneParentDownload()
	if err != nil {
		log.Fatal(err)
	}
	parentDownload.Parent = parent
	parentDownload.ParentID = parent.ID
	parentDownloadSample := []struct {
		id             string
		statusCode     int
		targetPath     string
		receivedBytes  int
		totalBytes     int
		siteUrl        string
		tabUrl         string
		tabReferredUrl string
		mimeType       string
		parent_id      int
		errorMessage   string
	}{
		{
			id:             strconv.Itoa(int(parentDownload.ID)),
			statusCode:     200,
			targetPath:     parentDownload.TargetPath,
			receivedBytes:  int(parentDownload.ReceivedBytes),
			totalBytes:     int(parentDownload.TotalBytes),
			siteUrl:        parentDownload.SiteUrl,
			tabUrl:         parentDownload.TabUrl,
			tabReferredUrl: parentDownload.TabReferredUrl,
			mimeType:       parentDownload.MimeType,
			parent_id:      int(parentDownload.ParentID),
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range parentDownloadSample {

		req, err := http.NewRequest("GET", "/parentdownloads", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetParentDownload)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, float64(parentDownload.ID), responseMap["id"])
			assert.Equal(t, parentDownload.TargetPath, responseMap["target_path"]) //the response url id is float64
			assert.Equal(t, float64(parentDownload.ReceivedBytes), responseMap["received_bytes"])
			assert.Equal(t, float64(parentDownload.TotalBytes), responseMap["total_bytes"])
			assert.Equal(t, parentDownload.SiteUrl, responseMap["site_url"])
			assert.Equal(t, parentDownload.TabUrl, responseMap["tab_url"])
			assert.Equal(t, parentDownload.TabReferredUrl, responseMap["tab_referred_url"])
			assert.Equal(t, parentDownload.MimeType, responseMap["mime_type"])
			assert.Equal(t, float64(parentDownload.ParentID), responseMap["parent_id"]) //the response parent id is float64
		}
	}
}

func TestDeleteParentDownload(t *testing.T) {
	var ParentEmail, ParentPassword string
	var ParentID uint32
	var AuthParentdownloadID uint64

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	parents, parentdownloads, err := seedParentsAndParentdownloadsAndUrls()
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

	// Get only the second parentdownload
	for _, parentdownload := range parentdownloads {
		if parentdownload.ID == 1 {
			continue
		}
		AuthParentdownloadID = parentdownload.ID
		ParentID = parentdownload.ParentID
	}
	parentdownloadSample := []struct {
		id           string
		parent_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthParentdownloadID)),
			parent_id:    ParentID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthParentdownloadID)),
			parent_id:    ParentID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthParentdownloadID)),
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
	for _, v := range parentdownloadSample {

		req, _ := http.NewRequest("GET", "/parentdownloads", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteParentDownload)

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
