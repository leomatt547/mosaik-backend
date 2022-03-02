package modeltests

import (
	"log"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	"github.com/go-playground/assert/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestFindAllParentDownloads(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table %v\n", err)
	}

	err = seedParents()
	if err != nil {
		log.Fatalf("Error seeding parent and parent table %v\n", err)
	}

	_, err = seedParentDownloads()
	if err != nil {
		log.Fatalf("Error seeding, url, and parent download table %v\n", err)
	}
	parentdownloads, err := parentDownloadInstance.FindAllParentDownloads(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the parent download: %v\n", err)
		return
	}
	assert.Equal(t, len(*parentdownloads), 2)
}

func TestSaveParentDownload(t *testing.T) {

	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}

	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Cannot seed parent %v\n", err)
	}

	newParentDownload := models.ParentDownload{
		TargetPath:     "D:/",
		ReceivedBytes:  100,
		TotalBytes:     100,
		SiteUrl:        "www.google.com",
		TabUrl:         "google.com/tabURL",
		TabReferredUrl: "google.com/tabURL",
		MimeType:       "text/html",
		ParentID:       parent.ID,
	}
	savedParentDownload, err := newParentDownload.SaveParentDownload(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the parent_download: %v\n", err)
		return
	}
	assert.Equal(t, newParentDownload.ID, savedParentDownload.ID)
	assert.Equal(t, newParentDownload.TargetPath, savedParentDownload.TargetPath)
	assert.Equal(t, newParentDownload.ReceivedBytes, savedParentDownload.ReceivedBytes)
	assert.Equal(t, newParentDownload.TotalBytes, savedParentDownload.TotalBytes)
	assert.Equal(t, newParentDownload.SiteUrl, savedParentDownload.SiteUrl)
	assert.Equal(t, newParentDownload.TabUrl, savedParentDownload.TabUrl)
	assert.Equal(t, newParentDownload.TabReferredUrl, savedParentDownload.TabReferredUrl)
	assert.Equal(t, newParentDownload.MimeType, savedParentDownload.MimeType)
	assert.Equal(t, newParentDownload.ParentID, savedParentDownload.ParentID)
}

func TestFindParentDownloadByID(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error all refreshing table %v\n", err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Error Seeding table parent")
	}
	parentdownload, err := seedOneParentDownload()
	if err != nil {
		log.Fatalf("Error Seeding table parent_download")
	}
	parentdownload.Parent = parent
	parentdownload.ParentID = parent.ID
	foundParentDownload, err := parentDownloadInstance.FindParentDownloadByID(server.DB, parentdownload.ID)
	if err != nil {
		t.Errorf("this is the error getting one parent Download: %v\n", err)
		return
	}
	assert.Equal(t, foundParentDownload.ID, parentdownload.ID)
	assert.Equal(t, foundParentDownload.TargetPath, parentdownload.TargetPath)
	assert.Equal(t, foundParentDownload.ReceivedBytes, parentdownload.ReceivedBytes)
	assert.Equal(t, foundParentDownload.TotalBytes, parentdownload.TotalBytes)
	assert.Equal(t, foundParentDownload.SiteUrl, parentdownload.SiteUrl)
	assert.Equal(t, foundParentDownload.TabUrl, parentdownload.TabUrl)
	assert.Equal(t, foundParentDownload.TabReferredUrl, parentdownload.TabReferredUrl)
	assert.Equal(t, foundParentDownload.MimeType, parentdownload.MimeType)
	assert.Equal(t, foundParentDownload.ParentID, parentdownload.ParentID)
}

func TestDeleteAParentDownload(t *testing.T) {
	err := refreshAllTable()
	if err != nil {
		log.Fatalf("Error refreshing all table: %v\n", err)
	}
	parent, err := seedOneParent()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	parentdownload, err := seedOneParentDownload()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	parentdownload.ParentID = parent.ID
	parentdownload.Parent = parent
	isDeleted, err := parentDownloadInstance.DeleteAParentDownload(server.DB, parentdownload.ID)
	if err != nil {
		t.Errorf("this is the error delete the parent download: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
