package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Url struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Url       string    `gorm:"type:text;not null;unique" json:"url"`
	Title     string    `gorm:"type:text;not null;unique" json:"title"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *Url) Validate() error {
	if u.Url == "" {
		return errors.New("butuh url")
	}
	if u.Title == "" {
		return errors.New("butuh title")
	}
	return nil
}

func (p *Url) SaveUrl(db *gorm.DB) (*Url, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Url{}, err
	}
	return p, nil
}

func (p *Url) FindAllUrls(db *gorm.DB) (*[]Url, error) {
	urls := []Url{}
	err := db.Debug().Model(&Url{}).Limit(100).Find(&urls).Error
	if err != nil {
		return &[]Url{}, err
	}
	return &urls, err
}

func (p *Url) FindUrlByID(db *gorm.DB, uid uint32) (*Url, error) {
	err := db.Debug().Model(Url{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Url{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Url{}, errors.New("Url Not Found")
	}
	return p, err
}
