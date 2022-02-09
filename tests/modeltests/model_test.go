package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.informatika.org/if3250_2022_37_alkademi/mosaik-backend/api/controllers"
	"gitlab.informatika.org/if3250_2022_37_alkademi/mosaik-backend/api/models"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var parentInstance = models.Parent{}
var childInstance = models.Child{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()
	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	// if TestDbDriver == "mysql" {
	// 	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbParent"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
	// 	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	// 	if err != nil {
	// 		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
	// 		log.Fatal("This is the error:", err)
	// 	} else {
	// 		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	// 	}
	// }
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

	refreshParentTable()

	parent := models.Parent{
		Nama: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.Parent{}).Create(&parent).Error
	if err != nil {
		log.Fatalf("cannot seed Parents table: %v", err)
	}
	return parent, nil
}

func seedParents() error {

	parents := []models.Parent{
		{
			Nama: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nama: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}

	for i, _ := range parents {
		err := server.DB.Model(&models.Parent{}).Create(&parents[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
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

func seedOneParentAndOneChild() (models.Child, error) {

	err := refreshParentAndChildTable()
	if err != nil {
		return models.Child{}, err
	}
	parent := models.Parent{
		Nama: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.Parent{}).Create(&parent).Error
	if err != nil {
		return models.Child{}, err
	}
	child := models.Child{
		Nama: "Sam Phil Jr",
		Email:    "sam_jr@gmail.com",
		Password: "password",
		ParentID: 1,
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
			Nama: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nama: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var childs = []models.Child{
		{
			Nama: "Steven victor Jr",
			Email:    "steven_jr@gmail.com",
			Password: "password",
			ParentID: 1,
		},
		{
			Nama: "Magu Frank Jr",
			Email:    "magu_jr@gmail.com",
			Password: "password",
			ParentID: 2,
		},
	}

	for i, _ := range parents {
		err = server.DB.Model(&models.Parent{}).Create(&parents[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Parents table: %v", err)
		}
		childs[i].ParentID = parents[i].ID

		err = server.DB.Model(&models.Child{}).Create(&childs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed childs table: %v", err)
		}
	}
	return parents, childs, nil
}