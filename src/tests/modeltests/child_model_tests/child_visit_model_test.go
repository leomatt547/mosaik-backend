package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestFindAllChildVisits(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table %v\n", err)
	}

	_, _, err = seedParentsAndChilds()
	if err != nil {
		log.Fatalf("Error seeding parent and child table %v\n", err)
	}

	_, _, err = seedChildVisitsAndUrls()
	if err != nil {
		log.Fatalf("Error seeding, url, and child visit table %v\n", err)
	}
	childvisits, err := childVisitInstance.FindAllChildVisits(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the child visit: %v\n", err)
		return
	}
	assert.Equal(t, len(*childvisits), 2)
}

func TestSaveChildVisit(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}

	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Cannot seed parent and child %v\n", err)
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatalf("Cannot seed url %v\n", err)
	}

	newChildVisit := models.ChildVisit{
		ID:       1,
		UrlID:    url.ID,
		Duration: 5,
		ChildID:  child.ID,
	}
	savedChildVisit, err := newChildVisit.SaveChildVisit(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the child_visit: %v\n", err)
		return
	}
	assert.Equal(t, newChildVisit.ID, savedChildVisit.ID)
	assert.Equal(t, newChildVisit.UrlID, savedChildVisit.UrlID)
	assert.Equal(t, newChildVisit.Duration, savedChildVisit.Duration)
	assert.Equal(t, newChildVisit.ChildID, savedChildVisit.ChildID)
}

func TestFindChildVisitByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding table child")
	}
	url, err := seedOneUrl()
	if err != nil {
		log.Fatalf("Error Seeding table url")
	}
	childvisit, err := seedOneChildVisit()
	if err != nil {
		log.Fatalf("Error Seeding table child_visit")
	}
	childvisit.Child = child
	childvisit.ChildID = child.ID
	childvisit.Url = url
	childvisit.UrlID = url.ID
	foundChildVisit, err := childVisitInstance.FindChildVisitByID(server.DB, childvisit.ID)
	if err != nil {
		t.Errorf("this is the error getting one child Visit: %v\n", err)
		return
	}
	assert.Equal(t, foundChildVisit.ID, childvisit.ID)
	assert.Equal(t, foundChildVisit.UrlID, childvisit.UrlID)
	assert.Equal(t, foundChildVisit.Duration, childvisit.Duration)
	assert.Equal(t, foundChildVisit.ChildID, childvisit.ChildID)
}

func TestDeleteAChildVisit(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	childvisit, err := seedOneChildVisit()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	childvisit.ChildID = child.ID
	childvisit.Child = child
	isDeleted, err := childVisitInstance.DeleteAChildVisit(server.DB, childvisit.ID)
	if err != nil {
		t.Errorf("this is the error delete the child visit: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
