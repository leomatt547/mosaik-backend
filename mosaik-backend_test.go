package main

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestDatabaseConnect(t *testing.T){

    db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s", "rdulpadeipmwvr", "2cf4c3b493be5216e5309be9837d987e6be5294c696314809eed1702a230d15b", "ec2-52-213-119-221.eu-west-1.compute.amazonaws.com", "d93bbe48ni70dp"))

	if err != nil {
        t.Errorf("failed to connect database")
	}

	defer db.Close()
}