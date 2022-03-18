package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type ParentVisit struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UrlID     uint64    `gorm:"size:255;not null;" json:"url_id"`
	Url       Url       `json:"Url"`
	Duration  uint64    `gorm:"type:bigint; not null;" json:"duration"`
	ParentID  uint32    `gorm:"not null" json:"parent_id"`
	Parent    Parent    `json:"Parent"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (pv *ParentVisit) Validate() error {
	if pv.Duration == 0 {
		return errors.New("butuh durasi")
	}
	if pv.ParentID == 0 {
		return errors.New("butuh parent_id")
	}
	if pv.UrlID == 0 {
		return errors.New("butuh url_id")
	}
	return nil
}

func (pv *ParentVisit) SaveParentVisit(db *gorm.DB) (*ParentVisit, error) {
	//POST variabel butuh: IDParent, IDUrl
	var err error
	err = db.Debug().Model(&ParentVisit{}).Create(&pv).Error
	if err != nil {
		return &ParentVisit{}, err
	}
	if pv.ID != 0 {
		//Dapatkan id Parent apakah ada atau tidak
		err = db.Debug().Model(&Parent{}).Where("id = ?", pv.ParentID).Take(&pv.Parent).Error
		if err != nil {
			return &ParentVisit{}, err
		}
		//Dapatkan id Url apakah ada atau tidak
		err = db.Debug().Model(&Url{}).Where("id = ?", pv.UrlID).Take(&pv.Url).Error
		if err != nil {
			return &ParentVisit{}, err
		}
	}
	return pv, nil
}

func (pv *ParentVisit) FindAllParentVisits(db *gorm.DB) (*[]ParentVisit, error) {
	var err error
	parentvisits := []ParentVisit{}
	err = db.Debug().Model(&ParentVisit{}).Limit(100).Find(&parentvisits).Error
	if err != nil {
		return &[]ParentVisit{}, err
	}
	if len(parentvisits) > 0 {
		//Dapatkan id Parent
		for i := range parentvisits {
			err := db.Debug().Model(&Parent{}).Where("id = ?", parentvisits[i].ParentID).Take(&parentvisits[i].Parent).Error
			if err != nil {
				return &[]ParentVisit{}, err
			}
			//Dapatkan id Url
			err = db.Debug().Model(&Url{}).Where("id = ?", parentvisits[i].UrlID).Take(&parentvisits[i].Url).Error
			if err != nil {
				return &[]ParentVisit{}, err
			}
		}
	}
	return &parentvisits, nil
}

func (pv *ParentVisit) FindParentVisitByID(db *gorm.DB, pid uint64) (*ParentVisit, error) {
	//POST variabel butuh: ID
	var err error
	err = db.Debug().Model(&ParentVisit{}).Where("id = ?", pid).Take(&pv).Error
	if err != nil {
		return &ParentVisit{}, err
	}
	//Apabila tidak ada error, lanjut
	if pv.ID != 0 {
		//Ambil Parent nya
		err = db.Debug().Model(&Parent{}).Where("id = ?", pv.ParentID).Take(&pv.Parent).Error
		if err != nil {
			return &ParentVisit{}, err
		}
		//Ambil Url yang dikunjungi
		err = db.Debug().Model(&Url{}).Where("id = ?", pv.UrlID).Take(&pv.Url).Error
		if err != nil {
			return &ParentVisit{}, err
		}
	}
	return pv, nil
}

func (pv *ParentVisit) FindParentVisitsbyParentID(db *gorm.DB, pid uint32) (*[]ParentVisit, error) {
	var err error
	parentvisits := []ParentVisit{}
	err = db.Debug().Model(&ParentVisit{}).Limit(100).Where("parent_id = ?", pid).Find(&parentvisits).Error
	if err != nil {
		return &[]ParentVisit{}, err
	}
	if len(parentvisits) > 0 {
		//Dapatkan id Parent
		for i := range parentvisits {
			err := db.Debug().Model(&Parent{}).Where("id = ?", pid).Take(&parentvisits[i].Parent).Error
			if err != nil {
				return &[]ParentVisit{}, err
			}
			//Dapatkan id Url
			err = db.Debug().Model(&Url{}).Where("id = ?", parentvisits[i].UrlID).Take(&parentvisits[i].Url).Error
			if err != nil {
				return &[]ParentVisit{}, err
			}
		}
	}
	return &parentvisits, nil
}

func (pv *ParentVisit) DeleteAParentVisit(db *gorm.DB, pvid uint64) (int64, error) {
	//Butuh id dan Parent id
	db = db.Debug().Model(&ParentVisit{}).Where("id = ?", pvid).Take(&ParentVisit{}).Delete(&ParentVisit{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Parent Visit not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
