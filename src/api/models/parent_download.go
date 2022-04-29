package models

import (
	"errors"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
)

type ParentDownload struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	TargetPath     string    `gorm:"type:text;not null;" json:"target_path"`
	ReceivedBytes  uint64    `gorm:"type:bigint; not null;" json:"received_bytes"`
	TotalBytes     uint64    `gorm:"type:bigint; not null;" json:"total_bytes"`
	SiteUrl        string    `gorm:"type:text;not null;" json:"site_url"`
	TabUrl         string    `gorm:"type:text;not null;" json:"tab_url"`
	TabReferredUrl string    `gorm:"type:text;" json:"tab_referred_url"`
	MimeType       string    `gorm:"type:varchar(255);not null;" json:"mime_type"`
	ParentID       uint32    `gorm:"not null" json:"parent_id"`
	Parent         Parent    `json:"Parent"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (pd *ParentDownload) Validate() error {
	if pd.TargetPath == "" {
		return errors.New("butuh target_path")
	}
	if pd.SiteUrl == "" {
		return errors.New("butuh site_url")
	}
	if pd.TabUrl == "" {
		return errors.New("butuh tab_url")
	}
	if pd.MimeType == "" {
		return errors.New("butuh mime_type")
	}
	if pd.ParentID == 0 {
		return errors.New("butuh parent")
	}
	return nil
}

func (pd *ParentDownload) SaveParentDownload(db *gorm.DB) (*ParentDownload, error) {
	//POST variabel butuh: IDParent
	var err error
	err = db.Debug().Model(&ParentDownload{}).Create(&pd).Error
	if err != nil {
		return &ParentDownload{}, err
	}
	if pd.ID != 0 {
		//Dapatkan id Parent apakah ada atau tidak
		err = db.Debug().Model(&Parent{}).Where("id = ?", pd.ParentID).Take(&pd.Parent).Error
		if err != nil {
			return &ParentDownload{}, err
		}
	}
	return pd, nil
}

func (pd *ParentDownload) FindAllParentDownloads(db *gorm.DB) (*[]ParentDownload, error) {
	var err error
	parentdownloads := []ParentDownload{}
	err = db.Debug().Model(&ParentDownload{}).Limit(100).Find(&parentdownloads).Error
	if err != nil {
		return &[]ParentDownload{}, err
	}
	if len(parentdownloads) > 0 {
		//Dapatkan id Parent
		for i := range parentdownloads {
			err := db.Debug().Model(&Parent{}).Where("id = ?", parentdownloads[i].ParentID).Take(&parentdownloads[i].Parent).Error
			if err != nil {
				return &[]ParentDownload{}, err
			}
		}
	}
	return &parentdownloads, nil
}

func (pd *ParentDownload) FindParentDownloadByID(db *gorm.DB, pid uint64) (*ParentDownload, error) {
	//POST variabel butuh: ID
	var err error
	err = db.Debug().Model(&ParentDownload{}).Where("id = ?", pid).Take(&pd).Error
	if err != nil {
		return &ParentDownload{}, err
	}
	//Apabila tidak ada error, lanjut
	if pd.ID != 0 {
		//Ambil Parent nya
		err = db.Debug().Model(&Parent{}).Where("id = ?", pd.ParentID).Take(&pd.Parent).Error
		if err != nil {
			return &ParentDownload{}, err
		}
	}
	return pd, nil
}

func (pd *ParentDownload) FindParentDownloadsbyParentID(db *gorm.DB, pid uint32) (*[]ParentDownload, error) {
	var err error
	parentdownloads := []ParentDownload{}
	err = db.Debug().Model(&ParentDownload{}).Limit(100).Where("parent_id = ?", pid).Find(&parentdownloads).Error
	if err != nil {
		return &[]ParentDownload{}, err
	}
	if len(parentdownloads) > 0 {
		//Dapatkan id Parent
		for i := range parentdownloads {
			err := db.Debug().Model(&Parent{}).Where("id = ?", pid).Take(&parentdownloads[i].Parent).Error
			if err != nil {
				return &[]ParentDownload{}, err
			}
		}
	}
	return &parentdownloads, nil
}

func (pd *ParentDownload) DeleteAParentDownload(db *gorm.DB, pdid uint64) (int64, error) {
	//Butuh id dan Parent id
	db = db.Debug().Model(&ParentDownload{}).Where("id = ?", pdid).Take(&ParentDownload{}).Delete(&ParentDownload{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Parent Download not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
