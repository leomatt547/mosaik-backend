package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_alkademi/mosaik-backend/api/models"

	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestFindAllChilds(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatalf("Error refreshing parent and child table %v\n", err)
	}
	_, _, err = seedParentsAndChilds()
	if err != nil {
		log.Fatalf("Error seeding parent and child table %v\n", err)
	}
	childs, err := childInstance.FindAllChilds(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the childs: %v\n", err)
		return
	}
	assert.Equal(t, len(*childs), 2)
}

func TestSaveChild(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatalf("Error parent and child refreshing table %v\n", err)
	}

	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Cannot seed parent %v\n", err)
	}

	newChild := models.Child{
		ID:       1,
		Nama:    "This is the nama",
		Email:  "This is the email",
		ParentID: parent.ID,
	}
	savedChild, err := newChild.SaveChild(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the child: %v\n", err)
		return
	}
	assert.Equal(t, newChild.ID, savedChild.ID)
	assert.Equal(t, newChild.Nama, savedChild.Nama)
	assert.Equal(t, newChild.Email, savedChild.Email)
	assert.Equal(t, newChild.ParentID, savedChild.ParentID)

}

func TestGetChildByID(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatalf("Error refreshing parent and child table: %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundChild, err := childInstance.FindChildByID(server.DB, child.ID)
	if err != nil {
		t.Errorf("this is the error getting one parent: %v\n", err)
		return
	}
	assert.Equal(t, foundChild.ID, child.ID)
	assert.Equal(t, foundChild.Nama, child.Nama)
	assert.Equal(t, foundChild.Email, child.Email)
}

func TestUpdateAChild(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatalf("Error refreshing parent and child table: %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	childUpdate := models.Child{
		ID:       1,
		Nama:    "modiUpdate",
		Email:  "modiupdate@gmail.com",
		ParentID: child.ParentID,
	}
	updatedChild, err := childUpdate.UpdateAChild(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the parent: %v\n", err)
		return
	}
	assert.Equal(t, updatedChild.ID, childUpdate.ID)
	assert.Equal(t, updatedChild.Nama, childUpdate.Nama)
	assert.Equal(t, updatedChild.Email, childUpdate.Email)
	assert.Equal(t, updatedChild.ParentID, childUpdate.ParentID)
}

func TestDeleteAChild(t *testing.T) {

	err := refreshParentAndChildTable()
	if err != nil {
		log.Fatalf("Error refreshing parent and child table: %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := childInstance.DeleteAChild(server.DB, child.ID, child.ParentID)
	if err != nil {
		t.Errorf("this is the error updating the parent: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}