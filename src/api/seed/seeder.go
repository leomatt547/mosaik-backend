package seed

import (
	"log"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"

	"github.com/jinzhu/gorm"
)

var parents = []models.Parent{
	{
		Nama:     "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Nama:     "Martin Luther",
		Email:    "luther@gmail.com",
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
		Nama:     "Martin Luther Jr",
		Email:    "luther_jr@gmail.com",
		Password: "password",
		ParentID: 2,
	},
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

var parentvisits = []models.ParentVisit{
	{
		UrlID:    2,
		Duration: 7,
		ParentID: 1,
	},
	{
		UrlID:    1,
		Duration: 11,
		ParentID: 2,
	},
}

var parentdownloads = []models.ParentDownload{
	{
		TargetPath:     "D:/",
		ReceivedBytes:  100,
		TotalBytes:     100,
		SiteUrl:        "www.google.com",
		TabUrl:         "google.com/tabURL",
		TabReferredUrl: "google.com/tabRefferedURL",
		MimeType:       "text/html",
		ParentID:       1,
	},
	{
		TargetPath:     "C:/",
		ReceivedBytes:  200,
		TotalBytes:     200,
		SiteUrl:        "www.facebook.com",
		TabUrl:         "facebook.com/tabURL",
		TabReferredUrl: "facebook.com/tabReferredURL",
		MimeType:       "text/html",
		ParentID:       2,
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Child{}, &models.Parent{}, &models.Url{}, &models.ChildVisit{}, &models.ParentVisit{}, &models.ParentDownload{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.Parent{}, &models.Child{}, &models.Url{}, &models.ChildVisit{}, &models.ParentVisit{}, &models.ParentDownload{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	//Adding Foreign Key Child
	err = db.Debug().Model(&models.Child{}).AddForeignKey("parent_id", "parents(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key parent error: %v", err)
	}

	//Foreign Key History
	err = db.Debug().Model(&models.ChildVisit{}).AddForeignKey("child_id", "children(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key child error: %v", err)
	}

	err = db.Debug().Model(&models.ChildVisit{}).AddForeignKey("url_id", "urls(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key url error: %v", err)
	}

	err = db.Debug().Model(&models.ParentVisit{}).AddForeignKey("parent_id", "parents(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key parent error: %v", err)
	}

	err = db.Debug().Model(&models.ParentVisit{}).AddForeignKey("url_id", "urls(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key url error: %v", err)
	}

	//Foreign Key Download
	err = db.Debug().Model(&models.ParentDownload{}).AddForeignKey("parent_id", "parents(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key parent error: %v", err)
	}

	for i, _ := range parents {
		//seeding parent
		err = db.Debug().Model(&models.Parent{}).Create(&parents[i]).Error
		if err != nil {
			log.Fatalf("cannot seed parents table: %v", err)
		}
		childs[i].ParentID = parents[i].ID

		//seeding child
		err = db.Debug().Model(&models.Child{}).Create(&childs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed childs table: %v", err)
		}

		//seeding parent download
		parentdownloads[i].ParentID = parents[i].ID
		err = db.Debug().Model(&models.ParentDownload{}).Create(&parentdownloads[i]).Error
		if err != nil {
			log.Fatalf("cannot seed parent_downloads table: %v", err)
		}
	}

	for i, _ := range urls {
		err = db.Debug().Model(&models.Url{}).Create(&urls[i]).Error
		if err != nil {
			log.Fatalf("cannot seed urls table: %v", err)
		}

		childvisits[i].UrlID = urls[i].ID
		err = db.Debug().Model(&models.ChildVisit{}).Create(&childvisits[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Child Visits table: %v", err)
		}

		parentvisits[i].UrlID = urls[i].ID
		err = db.Debug().Model(&models.ParentVisit{}).Create(&parentvisits[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Parent Visits table: %v", err)
		}
	}
}
