package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type ChildVisit struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UrlID     uint64    `gorm:"size:255;not null;" json:"url_id"`
	Url       Url       `json:"Url"`
	Duration  uint64    `gorm:"type:bigint; not null;" json:"duration"`
	ChildID   uint64    `gorm:"not null" json:"child_id"`
	Child     Child     `json:"Child"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (cv *ChildVisit) Validate() error {
	if cv.Duration == 0 {
		return errors.New("butuh durasi")
	}
	if cv.ChildID == 0 {
		return errors.New("butuh child_id")
	}
	if cv.UrlID == 0 {
		return errors.New("butuh url_id")
	}
	return nil
}

func (cv *ChildVisit) SaveChildVisit(db *gorm.DB) (*ChildVisit, error) {
	//POST variabel butuh: IDChild, IDUrl
	var err error
	err = db.Debug().Model(&ChildVisit{}).Create(&cv).Error
	if err != nil {
		return &ChildVisit{}, err
	}
	if cv.ID != 0 {
		//Dapatkan id Child apakah ada atau tidak
		err = db.Debug().Model(&Child{}).Where("id = ?", cv.ChildID).Take(&cv.Child).Error
		if err != nil {
			return &ChildVisit{}, err
		}
		//Dapatkan id Url apakah ada atau tidak
		err = db.Debug().Model(&Url{}).Where("id = ?", cv.UrlID).Take(&cv.Url).Error
		if err != nil {
			return &ChildVisit{}, err
		}
	}
	return cv, nil
}

func (cv *ChildVisit) FindAllChildVisits(db *gorm.DB) (*[]ChildVisit, error) {
	var err error
	childvisits := []ChildVisit{}
	err = db.Debug().Model(&ChildVisit{}).Limit(100).Find(&childvisits).Error
	if err != nil {
		return &[]ChildVisit{}, err
	}
	if len(childvisits) > 0 {
		//Dapatkan id Child
		for i, _ := range childvisits {
			err := db.Debug().Model(&Child{}).Where("id = ?", childvisits[i].ChildID).Take(&childvisits[i].Child).Error
			if err != nil {
				return &[]ChildVisit{}, err
			}
			//Dapatkan id Url
			err = db.Debug().Model(&Url{}).Where("id = ?", childvisits[i].UrlID).Take(&childvisits[i].Url).Error
			if err != nil {
				return &[]ChildVisit{}, err
			}
		}
	}
	return &childvisits, nil
}

func (cv *ChildVisit) FindChildVisitByID(db *gorm.DB, pid uint64) (*ChildVisit, error) {
	//POST variabel butuh: ID
	var err error
	err = db.Debug().Model(&ChildVisit{}).Where("id = ?", pid).Take(&cv).Error
	if err != nil {
		return &ChildVisit{}, err
	}
	//Apabila tidak ada error, lanjut
	if cv.ID != 0 {
		//Ambil Child nya
		err = db.Debug().Model(&Child{}).Where("id = ?", cv.ChildID).Take(&cv.Child).Error
		if err != nil {
			return &ChildVisit{}, err
		}
		//Ambil Url yang dikunjungi
		err = db.Debug().Model(&Url{}).Where("id = ?", cv.UrlID).Take(&cv.Url).Error
		if err != nil {
			return &ChildVisit{}, err
		}
	}
	return cv, nil
}

func (cv *ChildVisit) FindChildVisitsbyChildID(db *gorm.DB, cid uint64) (*[]ChildVisit, error) {
	var err error
	childvisits := []ChildVisit{}
	err = db.Debug().Model(&ChildVisit{}).Limit(100).Where("child_id = ?", cid).Find(&childvisits).Error
	if err != nil {
		return &[]ChildVisit{}, err
	}
	if len(childvisits) > 0 {
		//Dapatkan id Child
		for i, _ := range childvisits {
			err := db.Debug().Model(&Child{}).Where("id = ?", childvisits[i].ChildID).Take(&childvisits[i].Child).Error
			if err != nil {
				return &[]ChildVisit{}, err
			}
			//Dapatkan id Url
			for j, _ := range childvisits {
				err := db.Debug().Model(&Url{}).Where("id = ?", childvisits[j].UrlID).Take(&childvisits[j].Url).Error
				if err != nil {
					return &[]ChildVisit{}, err
				}
			}
		}
	}
	return &childvisits, nil
}

func (cv *ChildVisit) DeleteAChildVisit(db *gorm.DB, cvid uint64) (int64, error) {
	//Butuh id dan Child id
	db = db.Debug().Model(&ChildVisit{}).Where("id = ?", cvid).Take(&ChildVisit{}).Delete(&ChildVisit{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Child Visit not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
