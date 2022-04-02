package models

import (
	"errors"
	"time"
)

type Blacklist struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UrlID     uint64    `gorm:"size:255;not null;" json:"url_id"`
	Url       Url       `json:"Url"`
	ParentID  uint32    `gorm:"not null" json:"parent_id"`
	Parent    Parent    `json:"Parent"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (bl *Blacklist) Validate() error {
	if bl.UrlID == 0 {
		return errors.New("butuh url_id")
	}
	if bl.ParentID == 0 {
		return errors.New("butuh parent_id")
	}
	return nil
}

// func (bl *Blacklist)
