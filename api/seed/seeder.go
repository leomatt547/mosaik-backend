package seed

import (
	"log"

	"mosaik/api/models"

	"github.com/jinzhu/gorm"
)

var parents = []models.Parent{
	models.Parent{
		Nama: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.Parent{
		Nama: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var childs = []models.Child{
	models.Child{
		Nama: "Steven victor Jr",
		Email:    "steven_jr@gmail.com",
		Password: "password",
	},
	models.Child{
		Nama: "Steven victor Jr",
		Email:    "steven_jr@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Child{}, &models.Parent{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Parent{}, &models.Child{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Child{}).AddForeignKey("child_id", "parents(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range parents {
		err = db.Debug().Model(&models.Child{}).Create(&parents[i]).Error
		if err != nil {
			log.Fatalf("cannot seed parents table: %v", err)
		}
		childs[i].ParentID = parents[i].ID

		err = db.Debug().Model(&models.Parent{}).Create(&childs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed parents table: %v", err)
		}
	}
}