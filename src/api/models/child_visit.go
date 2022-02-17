package models

import (
	"errors"
	"time"
)

type ChildVisit struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UrlID     uint64    `gorm:"size:255;not null;unique" json:"url_id"`
	Url       Url       `json:"Url"`
	Duration  uint64    `gorm:"type:bigint; not null;" json:"duration"`
	ChildID   uint32    `gorm:"not null" json:"child_id"`
	Child     Child     `json:"Child"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (cv *ChildVisit) Validate() error {
	if cv.Duration == 0 {
		return errors.New("butuh durasi")
	}
	return nil
}
