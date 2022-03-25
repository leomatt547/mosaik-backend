package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	//_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
)

func TestFindAllNSFWUrls(t *testing.T) {

	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedNSFWUrls()
	if err != nil {
		log.Fatal(err)
	}

	nsfw_urls, err := nsfwUrlInstance.FindAllNSFWUrls(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the nsfw_urls: %v\n", err)
		return
	}
	assert.Equal(t, len(*nsfw_urls), 2)
}

func TestSaveNSFWUrl(t *testing.T) {

	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	newNSFWUrl := models.NSFWUrl{
		ID:  1,
		Url: "www.pornhub.com",
	}
	savedNSFWUrl, err := newNSFWUrl.SaveNSFWUrl(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the nsfw_urls: %v\n", err)
		return
	}
	assert.Equal(t, newNSFWUrl.ID, savedNSFWUrl.ID)
	assert.Equal(t, newNSFWUrl.Url, savedNSFWUrl.Url)
}

func TestFindNSFWUrlByID(t *testing.T) {

	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}

	new_nsfw_url, err := seedOneNSFWUrl()
	if err != nil {
		log.Fatalf("cannot seed nsfw_url table: %v", err)
	}
	foundnsfw_url, err := nsfwUrlInstance.FindNSFWUrlByID(server.DB, new_nsfw_url.ID)
	if err != nil {
		t.Errorf("this is the error getting one nsfw_url: %v\n", err)
		return
	}
	assert.Equal(t, foundnsfw_url.ID, new_nsfw_url.ID)
	assert.Equal(t, foundnsfw_url.Url, new_nsfw_url.Url)
}

func TestFindRecordByNSFWUrl(t *testing.T) {

	err := refreshNSFWUrlTable()
	if err != nil {
		log.Fatal(err)
	}

	new_nsfw_url, err := seedOneNSFWUrl()
	if err != nil {
		log.Fatalf("Cannot seed nsfw_url: %v\n", err)
	}

	query := "www.google.com"
	hasil_nsfw_url, err := nsfwUrlInstance.FindRecordByNSFWUrl(server.DB, query)
	if err != nil {
		t.Errorf("this is the error find the nsfw_url: %v\n", err)
		return
	}
	assert.Equal(t, new_nsfw_url.ID, hasil_nsfw_url.ID)
	assert.Equal(t, new_nsfw_url.Url, hasil_nsfw_url.Url)
}
