package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	//_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
)

func TestFindAllUrls(t *testing.T) {

	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUrls()
	if err != nil {
		log.Fatal(err)
	}

	urls, err := urlInstance.FindAllUrls(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the urls: %v\n", err)
		return
	}
	assert.Equal(t, len(*urls), 2)
}

func TestSaveUrl(t *testing.T) {

	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	newUrl := models.Url{
		ID:    1,
		Url:   "www.youtube.com",
		Title: "Youtube",
	}
	savedUrl, err := newUrl.SaveUrl(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the urls: %v\n", err)
		return
	}
	assert.Equal(t, newUrl.ID, savedUrl.ID)
	assert.Equal(t, newUrl.Url, savedUrl.Url)
	assert.Equal(t, newUrl.Title, savedUrl.Title)
}

func TestFindUrlByID(t *testing.T) {

	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}

	new_url, err := seedOneUrl()
	if err != nil {
		log.Fatalf("cannot seed url table: %v", err)
	}
	foundurl, err := urlInstance.FindUrlByID(server.DB, new_url.ID)
	if err != nil {
		t.Errorf("this is the error getting one url: %v\n", err)
		return
	}
	assert.Equal(t, foundurl.ID, new_url.ID)
	assert.Equal(t, foundurl.Url, new_url.Url)
	assert.Equal(t, foundurl.Title, new_url.Title)
}

func TestFindRecordByUrl(t *testing.T) {

	err := refreshUrlTable()
	if err != nil {
		log.Fatal(err)
	}

	new_url, err := seedOneUrl()
	if err != nil {
		log.Fatalf("Cannot seed url: %v\n", err)
	}

	query := "www.google.com"
	hasil_url, err := urlInstance.FindRecordByUrl(server.DB, query)
	if err != nil {
		t.Errorf("this is the error find the url: %v\n", err)
		return
	}
	assert.Equal(t, new_url.ID, hasil_url.ID)
	assert.Equal(t, new_url.Url, hasil_url.Url)
	assert.Equal(t, new_url.Title, hasil_url.Title)
}
