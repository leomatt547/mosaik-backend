package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSendMail(t *testing.T) {

	refreshParentTable()
	_, err := seedOneParent()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "pet@gmail.com"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "kangmail.com"}`,
			statusCode:   422,
			errorMessage: "invalid email",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/parents/resetpassword", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.SendMail)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}
		if v.statusCode == 422 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
