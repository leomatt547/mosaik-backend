package modeltests

import (
	"log"
	"testing"

	"mosaik-backend/api/models"

	//_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllParents(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedParents()
	if err != nil {
		log.Fatal(err)
	}

	parents, err := parentInstance.FindAllParents(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the parents: %v\n", err)
		return
	}
	assert.Equal(t, len(*parents), 2)
}

func TestSaveParent(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}
	newParent := models.Parent{
		ID:       1,
		Email:    "test@gmail.com",
		Nama: "test",
		Password: "password",
	}
	savedParent, err := newParent.SaveParent(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the parents: %v\n", err)
		return
	}
	assert.Equal(t, newParent.ID, savedParent.ID)
	assert.Equal(t, newParent.Email, savedParent.Email)
	assert.Equal(t, newParent.Nama, savedParent.Nama)
}

func TestGetParentByID(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}

	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("cannot seed parents table: %v", err)
	}
	foundparent, err := parentInstance.FindParentByID(server.DB, parent.ID)
	if err != nil {
		t.Errorf("this is the error getting one parent: %v\n", err)
		return
	}
	assert.Equal(t, foundparent.ID, parent.ID)
	assert.Equal(t, foundparent.Email, parent.Email)
	assert.Equal(t, foundparent.Nama, parent.Nama)
}

func TestUpdateAParent(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}

	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Cannot seed parent: %v\n", err)
	}

	parentUpdate := models.Parent{
		ID:       1,
		Nama: "modiUpdate",
		Email:    "modiupdate@gmail.com",
		Password: "password",
	}
	updatedParent, err := parentUpdate.UpdateAParent(server.DB, parent.ID)
	if err != nil {
		t.Errorf("this is the error updating the parent: %v\n", err)
		return
	}
	assert.Equal(t, updatedParent.ID, parentUpdate.ID)
	assert.Equal(t, updatedParent.Email, parentUpdate.Email)
	assert.Equal(t, updatedParent.Nama, parentUpdate.Nama)
}

func TestDeleteAParent(t *testing.T) {

	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}

	parent, err := seedOneParent()

	if err != nil {
		log.Fatalf("Cannot seed parent: %v\n", err)
	}

	isDeleted, err := parentInstance.DeleteAParent(server.DB, parent.ID)
	if err != nil {
		t.Errorf("this is the error updating the parent: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}