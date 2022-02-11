package seed

import (
	"log"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/models"

	"github.com/jinzhu/gorm"
)

var parents = []models.Parent{
	{
		Nama: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Nama: "Martin Luther",
		Email:    "luther@gmail.com",
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
		Nama: "Martin Luther Jr",
		Email:    "luther_jr@gmail.com",
		Password: "password",
		ParentID: 2,
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

	err = db.Debug().Model(&models.Child{}).AddForeignKey("parent_id", "parents(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range parents {
		err = db.Debug().Model(&models.Parent{}).Create(&parents[i]).Error
		if err != nil {
			log.Fatalf("cannot seed parents table: %v", err)
		}
		childs[i].ParentID = parents[i].ID

		err = db.Debug().Model(&models.Child{}).Create(&childs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed childs table: %v", err)
		}
	}
}