package models

import (
	"errors"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
)

type ChildDownload struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	TargetPath     string    `gorm:"type:text;not null;" json:"target_path"`
	ReceivedBytes  uint64    `gorm:"type:bigint; not null;" json:"received_bytes"`
	TotalBytes     uint64    `gorm:"type:bigint; not null;" json:"total_bytes"`
	SiteUrl        string    `gorm:"type:text;not null;" json:"site_url"`
	TabUrl         string    `gorm:"type:text;not null;" json:"tab_url"`
	TabReferredUrl string    `gorm:"type:text;" json:"tab_referred_url"`
	MimeType       string    `gorm:"type:varchar(255);not null;" json:"mime_type"`
	ChildID        uint64    `gorm:"not null" json:"child_id"`
	Child          Child     `json:"Child"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (cd *ChildDownload) Validate() error {
	if cd.TargetPath == "" {
		return errors.New("butuh target_path")
	}
	if cd.SiteUrl == "" {
		return errors.New("butuh site_url")
	}
	if cd.TabUrl == "" {
		return errors.New("butuh tab_url")
	}
	if cd.MimeType == "" {
		return errors.New("butuh mime_type")
	}
	if cd.ChildID == 0 {
		return errors.New("butuh child")
	}
	return nil
}

func (cd *ChildDownload) SaveChildDownload(db *gorm.DB) (*ChildDownload, error) {
	//POST variabel butuh: IDChild, IDUrl
	var err error
	err = db.Debug().Model(&ChildDownload{}).Create(&cd).Error
	if err != nil {
		return &ChildDownload{}, err
	}
	if cd.ID != 0 {
		//Dapatkan id Child apakah ada atau tidak
		err = db.Debug().Model(&Child{}).Where("id = ?", cd.ChildID).Take(&cd.Child).Error
		if err != nil {
			return &ChildDownload{}, err
		}
	}
	return cd, nil
}

func (cd *ChildDownload) FindAllChildDownloads(db *gorm.DB) (*[]ChildDownload, error) {
	var err error
	childdownloads := []ChildDownload{}
	err = db.Debug().Model(&ChildDownload{}).Limit(100).Find(&childdownloads).Error
	if err != nil {
		return &[]ChildDownload{}, err
	}
	if len(childdownloads) > 0 {
		//Dapatkan id Child
		for i := range childdownloads {
			err := db.Debug().Model(&Child{}).Where("id = ?", childdownloads[i].ChildID).Take(&childdownloads[i].Child).Error
			if err != nil {
				return &[]ChildDownload{}, err
			}
			err = db.Debug().Model(&Parent{}).Where("id = ?", childdownloads[i].Child.ParentID).Take(&childdownloads[i].Child.Parent).Error
			if err != nil {
				return &[]ChildDownload{}, err
			}
		}
	}
	return &childdownloads, nil
}

func (cd *ChildDownload) FindChildDownloadByID(db *gorm.DB, pid uint64) (*ChildDownload, error) {
	//POST variabel butuh: ID
	var err error
	err = db.Debug().Model(&ChildDownload{}).Where("id = ?", pid).Take(&cd).Error
	if err != nil {
		return &ChildDownload{}, err
	}
	//Apabila tidak ada error, lanjut
	if cd.ID != 0 {
		//Ambil Child nya
		err = db.Debug().Model(&Child{}).Where("id = ?", cd.ChildID).Take(&cd.Child).Error
		if err != nil {
			return &ChildDownload{}, err
		}
	}
	return cd, nil
}

func (cd *ChildDownload) FindChildDownloadsbyChildID(db *gorm.DB, cid uint64) (*[]ChildDownload, error) {
	var err error
	childdownloads := []ChildDownload{}
	err = db.Debug().Model(&ChildDownload{}).Limit(100).Where("child_id = ?", cid).Find(&childdownloads).Error
	if err != nil {
		return &[]ChildDownload{}, err
	}
	if len(childdownloads) > 0 {
		//Dapatkan id Child
		for i := range childdownloads {
			err := db.Debug().Model(&Child{}).Where("id = ?", cid).Take(&childdownloads[i].Child).Error
			if err != nil {
				return &[]ChildDownload{}, err
			}
		}
	}
	return &childdownloads, nil
}

func (cd *ChildDownload) DeleteAChildDownload(db *gorm.DB, cdid uint64) (int64, error) {
	//Butuh id dan Child id
	db = db.Debug().Model(&ChildDownload{}).Where("id = ?", cdid).Take(&ChildDownload{}).Delete(&ChildDownload{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Child Download not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
