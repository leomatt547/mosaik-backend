package models

import (
	"errors"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
)

type NSFWUrl struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Url       string    `gorm:"type:text;not null;unique" json:"url"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *NSFWUrl) Validate() error {
	if u.Url == "" {
		return errors.New("butuh url")
	}
	return nil
}

func (u *NSFWUrl) SaveNSFWUrl(db *gorm.DB) (*NSFWUrl, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &NSFWUrl{}, err
	}
	return u, nil
}

func (u *NSFWUrl) FindAllNSFWUrls(db *gorm.DB) (*[]NSFWUrl, error) {
	urls := []NSFWUrl{}
	err := db.Debug().Model(&NSFWUrl{}).Limit(100).Find(&urls).Error
	if err != nil {
		return &[]NSFWUrl{}, err
	}
	return &urls, err
}

func (u *NSFWUrl) FindNSFWUrlByID(db *gorm.DB, uid uint64) (*NSFWUrl, error) {
	err := db.Debug().Model(NSFWUrl{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &NSFWUrl{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &NSFWUrl{}, errors.New("NSFWUrl Not Found")
	}
	return u, err
}

func (u *NSFWUrl) FindRecordByNSFWUrl(db *gorm.DB, link string) (*NSFWUrl, error) {
	err := db.Debug().Model(NSFWUrl{}).Where("url = ?", link).Take(&u).Error
	if err != nil {
		return &NSFWUrl{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &NSFWUrl{}, errors.New("NSFWUrl Not Found")
	}
	return u, err
}
