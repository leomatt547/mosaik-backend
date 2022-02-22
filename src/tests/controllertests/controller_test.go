package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/controllers"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
)

var server = controllers.Server{}

// var parentInstance = models.Parent{}
// var childInstance = models.Child{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()
	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TestDbDriver")
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshParentTable() error {
	err := server.DB.DropTableIfExists(&models.Parent{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Parent{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneParent() (models.Parent, error) {
	err := refreshParentTable()
	if err != nil {
		log.Fatal(err)
	}

	parent := models.Parent{
		Nama:     "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.Parent{}).Create(&parent).Error
	if err != nil {
		return models.Parent{}, err
	}
	return parent, nil
}

func seedParents() ([]models.Parent, error) {
	var err error
	if err != nil {
		return nil, err
	}
	parents := []models.Parent{
		{
			Nama:     "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nama:     "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}
	for i, _ := range parents {
		err := server.DB.Model(&models.Parent{}).Create(&parents[i]).Error
		if err != nil {
			return []models.Parent{}, err
		}
	}
	return parents, nil
}

func refreshParentAndChildTable() error {
	err := server.DB.DropTableIfExists(&models.Parent{}, &models.Child{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Parent{}, &models.Child{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func refreshUrlTable() error {
	err := server.DB.DropTableIfExists(&models.Url{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Url{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed url tables")
	return nil
}

func refreshAllTable() error {
	err := server.DB.DropTableIfExists(&models.Parent{}, &models.Child{}, &models.ChildVisit{}, &models.Url{}, &models.ParentVisit{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Parent{}, &models.Child{}, &models.ChildVisit{}, &models.Url{}, &models.ParentVisit{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed all tables")
	return nil
}

func seedOneUrl() (models.Url, error) {
	url := models.Url{
		Url:   "www.google.com",
		Title: "google",
	}

	err := server.DB.Model(&models.Url{}).Create(&url).Error
	if err != nil {
		log.Fatalf("cannot seed Urls table: %v", err)
	}
	return url, nil
}

func seedUrls() error {
	var urls = []models.Url{
		{
			Url:   "www.google.com",
			Title: "Google",
		},
		{
			Url:   "www.facebook.com",
			Title: "Facebook",
		},
	}

	for i, _ := range urls {
		err := server.DB.Model(&models.Url{}).Create(&urls[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func seedChildVisitsAndUrls() ([]models.ChildVisit, []models.Url, error) {
	var err error
	if err != nil {
		return []models.ChildVisit{}, []models.Url{}, err
	}
	var urls = []models.Url{
		{
			Url:   "www.google.com",
			Title: "Google",
		},
		{
			Url:   "www.facebook.com",
			Title: "Facebook",
		},
	}
	var childvisits = []models.ChildVisit{
		{
			UrlID:    1,
			Duration: 5,
			ChildID:  1,
		},
		{
			UrlID:    2,
			Duration: 10,
			ChildID:  2,
		},
	}
	for i, _ := range urls {
		err = server.DB.Model(&models.Url{}).Create(&urls[i]).Error
		if err != nil {
			log.Fatalf("cannot seed urls table: %v", err)
		}

		childvisits[i].UrlID = urls[i].ID
		err = server.DB.Model(&models.ChildVisit{}).Create(&childvisits[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Child Visits table: %v", err)
		}
	}
	return childvisits, urls, nil
}

func seedParentVisitsAndUrls() ([]models.ParentVisit, []models.Url, error) {
	var err error
	if err != nil {
		return []models.ParentVisit{}, []models.Url{}, err
	}
	var urls = []models.Url{
		{
			Url:   "www.google.com",
			Title: "Google",
		},
		{
			Url:   "www.facebook.com",
			Title: "Facebook",
		},
	}
	var parentvisits = []models.ParentVisit{
		{
			UrlID:    1,
			Duration: 5,
			ParentID: 1,
		},
		{
			UrlID:    2,
			Duration: 10,
			ParentID: 2,
		},
	}
	for i, _ := range urls {
		err = server.DB.Model(&models.Url{}).Create(&urls[i]).Error
		if err != nil {
			log.Fatalf("cannot seed urls table: %v", err)
		}

		parentvisits[i].UrlID = urls[i].ID
		err = server.DB.Model(&models.ParentVisit{}).Create(&parentvisits[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Parent Visits table: %v", err)
		}
	}
	return parentvisits, urls, nil
}

func seedOneChildVisit() (models.ChildVisit, error) {
	var childvisits = models.ChildVisit{
		ID:       1,
		UrlID:    1,
		Duration: 5,
		ChildID:  1,
	}

	err := server.DB.Model(&models.ChildVisit{}).Create(&childvisits).Error
	if err != nil {
		log.Fatalf("cannot seed Child Visit table: %v", err)
	}
	return childvisits, nil
}

func seedOneParentVisit() (models.ParentVisit, error) {
	var parentvisits = models.ParentVisit{
		ID:       1,
		UrlID:    1,
		Duration: 5,
		ParentID: 1,
	}

	err := server.DB.Model(&models.ParentVisit{}).Create(&parentvisits).Error
	if err != nil {
		log.Fatalf("cannot seed Parent Visit table: %v", err)
	}
	return parentvisits, nil
}

func seedOneParentAndOneChild() (models.Child, error) {
	err := refreshParentAndChildTable()
	if err != nil {
		return models.Child{}, err
	}
	parent := models.Parent{
		Nama:     "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.Parent{}).Create(&parent).Error
	if err != nil {
		return models.Child{}, err
	}
	child := models.Child{
		Nama:     "This is the nama sam",
		Email:    "This is the email sam",
		ParentID: parent.ID,
	}
	err = server.DB.Model(&models.Child{}).Create(&child).Error
	if err != nil {
		return models.Child{}, err
	}
	return child, nil
}

func seedParentsAndChilds() ([]models.Parent, []models.Child, error) {
	var err error
	if err != nil {
		return []models.Parent{}, []models.Child{}, err
	}
	var parents = []models.Parent{
		{
			Nama:     "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nama:     "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var childs = []models.Child{
		{
			Nama:     "Steven victor Jr",
			Email:    "steven_jr@gmail.com",
			Password: "password",
			ParentID: 1,
		},
		{
			Nama:     "Magu Frank Jr",
			Email:    "magu_jr@gmail.com",
			Password: "password",
			ParentID: 2,
		},
	}

	for i, _ := range parents {
		err = server.DB.Model(&models.Parent{}).Create(&parents[i]).Error
		if err != nil {
			log.Fatalf("cannot seed parents table: %v", err)
		}
		childs[i].ParentID = parents[i].ID

		err = server.DB.Model(&models.Child{}).Create(&childs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed childs table: %v", err)
		}
	}
	return parents, childs, nil
}
