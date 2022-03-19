package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestFindAllChildDownloads(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table %v\n", err)
	}

	_, _, err = seedParentsAndChilds()
	if err != nil {
		log.Fatalf("Error seeding parent and child table %v\n", err)
	}

	_, err = seedChildDownloads()
	if err != nil {
		log.Fatalf("Error seeding child download table %v\n", err)
	}
	childdownloads, err := childDownloadInstance.FindAllChildDownloads(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the child download: %v\n", err)
		return
	}
	assert.Equal(t, len(*childdownloads), 2)
}

func TestSaveChildDownload(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}

	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Cannot seed parent and child %v\n", err)
	}

	newChildDownload := models.ChildDownload{
		TargetPath:     "D:/",
		ReceivedBytes:  100,
		TotalBytes:     100,
		SiteUrl:        "www.google.com",
		TabUrl:         "google.com/tabURL",
		TabReferredUrl: "google.com/tabURL",
		MimeType:       "text/html",
		ChildID:        child.ID,
	}
	savedChildDownload, err := newChildDownload.SaveChildDownload(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the child_download: %v\n", err)
		return
	}
	assert.Equal(t, newChildDownload.ID, savedChildDownload.ID)
	assert.Equal(t, newChildDownload.TargetPath, savedChildDownload.TargetPath)
	assert.Equal(t, newChildDownload.ReceivedBytes, savedChildDownload.ReceivedBytes)
	assert.Equal(t, newChildDownload.TotalBytes, savedChildDownload.TotalBytes)
	assert.Equal(t, newChildDownload.SiteUrl, savedChildDownload.SiteUrl)
	assert.Equal(t, newChildDownload.TabUrl, savedChildDownload.TabUrl)
	assert.Equal(t, newChildDownload.TabReferredUrl, savedChildDownload.TabReferredUrl)
	assert.Equal(t, newChildDownload.MimeType, savedChildDownload.MimeType)
	assert.Equal(t, newChildDownload.ChildID, savedChildDownload.ChildID)
}

func TestFindChildDownloadByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding table child")
	}
	childdownload, err := seedOneChildDownload()
	if err != nil {
		log.Fatalf("Error Seeding table child_download")
	}
	childdownload.Child = child
	childdownload.ChildID = child.ID
	foundChildDownload, err := childDownloadInstance.FindChildDownloadByID(server.DB, childdownload.ID)
	if err != nil {
		t.Errorf("this is the error getting one child Download: %v\n", err)
		return
	}
	assert.Equal(t, foundChildDownload.ID, childdownload.ID)
	assert.Equal(t, foundChildDownload.TargetPath, childdownload.TargetPath)
	assert.Equal(t, foundChildDownload.ReceivedBytes, childdownload.ReceivedBytes)
	assert.Equal(t, foundChildDownload.TotalBytes, childdownload.TotalBytes)
	assert.Equal(t, foundChildDownload.SiteUrl, childdownload.SiteUrl)
	assert.Equal(t, foundChildDownload.TabUrl, childdownload.TabUrl)
	assert.Equal(t, foundChildDownload.TabReferredUrl, childdownload.TabReferredUrl)
	assert.Equal(t, foundChildDownload.MimeType, childdownload.MimeType)
	assert.Equal(t, foundChildDownload.ChildID, childdownload.ChildID)
}

func TestDeleteAChildDownload(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	child, err := seedOneParentAndOneChild()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	childdownload, err := seedOneChildDownload()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	childdownload.ChildID = child.ID
	childdownload.Child = child
	isDeleted, err := childDownloadInstance.DeleteAChildDownload(server.DB, childdownload.ID)
	if err != nil {
		t.Errorf("this is the error delete the child download: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
