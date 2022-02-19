package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestFindAllParentVisits(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table %v\n", err)
	}

	err = seedParents()
	if err != nil {
		log.Fatalf("Error seeding parent and parent table %v\n", err)
	}

	_, _, err = seedParentVisitsAndUrls()
	if err != nil {
		log.Fatalf("Error seeding, url, and parent visit table %v\n", err)
	}
	parentvisits, err := parentVisitInstance.FindAllParentVisits(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the parent visit: %v\n", err)
		return
	}
	assert.Equal(t, len(*parentvisits), 2)
}

func TestSaveParentVisit(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}

	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Cannot seed parent and parent %v\n", err)
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatalf("Cannot seed url %v\n", err)
	}

	newParentVisit := models.ParentVisit{
		ID:       1,
		UrlID:    url.ID,
		Duration: 5,
		ParentID: parent.ID,
	}
	savedParentVisit, err := newParentVisit.SaveParentVisit(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the parent_visit: %v\n", err)
		return
	}
	assert.Equal(t, newParentVisit.ID, savedParentVisit.ID)
	assert.Equal(t, newParentVisit.UrlID, savedParentVisit.UrlID)
	assert.Equal(t, newParentVisit.Duration, savedParentVisit.Duration)
	assert.Equal(t, newParentVisit.ParentID, savedParentVisit.ParentID)
}

func TestFindParentVisitByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Error Seeding table parent")
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatalf("Error Seeding table url")
	}
	parentvisit, err := seedOneParentVisit()
	if err != nil {
		log.Fatalf("Error Seeding table parent_visit")
	}
	parentvisit.Parent = parent
	parentvisit.ParentID = parent.ID
	parentvisit.Url = url
	parentvisit.UrlID = url.ID
	foundParentVisit, err := parentVisitInstance.FindParentVisitByID(server.DB, parentvisit.ID)
	if err != nil {
		t.Errorf("this is the error getting one parent Visit: %v\n", err)
		return
	}
	assert.Equal(t, foundParentVisit.ID, parentvisit.ID)
	assert.Equal(t, foundParentVisit.UrlID, parentvisit.UrlID)
	assert.Equal(t, foundParentVisit.Duration, parentvisit.Duration)
	assert.Equal(t, foundParentVisit.ParentID, parentvisit.ParentID)
}

func TestDeleteAParentVisit(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	parentvisit, err := seedOneParentVisit()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	parentvisit.ParentID = parent.ID
	parentvisit.Parent = parent
	isDeleted, err := parentVisitInstance.DeleteAParentVisit(server.DB, parentvisit.ID)
	if err != nil {
		t.Errorf("this is the error delete the parent visit: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
