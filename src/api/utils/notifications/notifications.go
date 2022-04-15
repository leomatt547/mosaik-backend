package notifications

import (
	"net/http"
	"os"
	"strings"
)

func SendPushNotification(title string, body string, token string) error {
	url := "https://fcm.googleapis.com/fcm/send"
	method := "POST"

	payload := strings.NewReader(`{` + "" + ` "to" : ` + token + `,"data" : {"title":` + title + `,"body" : ` + body + `}` + "" + `}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "key="+os.Getenv("FIREBASE_KEY"))
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
