package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

type Parent struct {
	gorm.Model
	Nama    	string
	Email		string
	Password	string
	Childs  	[]Child
}
  
type Child struct {
	gorm.Model
	Nama    	string
	Email		string
	Password	string
	ParentID	int
}

var db *gorm.DB

var err error

var (
    parents = []Parent{
        {Nama: "Jimmy Johnson", Email: "Jimmy@gmail.com", Password:"pplkel37"},
        {Nama: "Howard Hills",  Email: "Howard@gmail.com", Password:"pplkel37"},
        {Nama: "Craig Colbin",  Email: "Craig@gmail.com", Password:"pplkel37"},
    }

    childs = []Child{
        {Nama: "Jimmy Junior", Email: "Jimmy_jr@gmail.com", Password:"pplkel37", ParentID:1},
        {Nama: "Howard Junior",  Email: "Howard_jr@gmail.com", Password:"pplkel37", ParentID:2},
        {Nama: "Craig Junior",  Email: "Craig_jr@gmail.com", Password:"pplkel37", ParentID:3},

    }
)

func main() {	
	router := mux.NewRouter()
	db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s", "rdulpadeipmwvr", "2cf4c3b493be5216e5309be9837d987e6be5294c696314809eed1702a230d15b", "ec2-52-213-119-221.eu-west-1.compute.amazonaws.com", "d93bbe48ni70dp"))

	if err != nil {
	  panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Parent{})
	db.AutoMigrate(&Child{})

	for index := range parents {
		db.Create(&parents[index])
	}

	for index := range childs {
		db.Create(&childs[index])
	}

	router.HandleFunc("/childs", GetChilds).Methods("GET")
	router.HandleFunc("/childs/{id}", GetChild).Methods("GET")
	router.HandleFunc("/parents/{id}", GetParent).Methods("GET")
	router.HandleFunc("/childs/{id}", DeleteChild).Methods("DELETE")

	handler := cors.Default().Handler(router)
	
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func GetChilds(w http.ResponseWriter, r *http.Request) {	
	var childs []Child
	db.Find(&childs)
	json.NewEncoder(w).Encode(&childs)
}

func GetChild(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var child Child
	db.First(&child, params["id"])
	json.NewEncoder(w).Encode(&child)
}

func GetParent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var parent	Parent
	var childs	[]Child
	db.First(&parent, params["id"])
	db.Model(&parent).Related(&childs)
	parent.Childs = childs
	json.NewEncoder(w).Encode(&parent)
}

func DeleteChild(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var child Child
	db.First(&child, params["id"])
	db.Delete(&child)

	var childs []Child
	db.Find(&childs)
	json.NewEncoder(w).Encode(&childs)
}