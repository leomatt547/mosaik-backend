package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Blacklist struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Url       string    `gorm:"uniqueIndex:blacklists_url_parent_id;type:text;not null;" json:"url"`
	ParentID  uint32    `gorm:"uniqueIndex:blacklists_url_parent_id;not null" json:"parent_id"`
	Parent    Parent    `json:"Parent"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (bl *Blacklist) Validate() error {
	if bl.Url == "" {
		return errors.New("butuh url")
	}
	if bl.ParentID == 0 {
		return errors.New("butuh parent_id")
	}
	return nil
}

func (bl *Blacklist) SaveBlacklist(db *gorm.DB) (*Blacklist, error) {
	err := db.Debug().Create(&bl).Error
	if err != nil {
		return &Blacklist{}, err
	}
	if bl.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", &bl.ParentID).Take(&bl.Parent).Error
		if err != nil {
			return &Blacklist{}, err
		}
	}
	return bl, err
}

func (bl *Blacklist) FindAllBlacklist(db *gorm.DB) (*[]Blacklist, error) {
	list := []Blacklist{}
	err := db.Debug().Model(&Blacklist{}).Limit(100).Find(&bl).Error
	if err != nil {
		return &[]Blacklist{}, err
	}
	if len(list) > 0 {
		for i := range list {
			err := db.Debug().Model(&Parent{}).Where("id = ?", list[i].ParentID).Take(&list[i].Parent).Error
			if err != nil {
				return &[]Blacklist{}, err
			}
		}
	}
	return &list, err
}

func (bl *Blacklist) FindBlacklistByID(db *gorm.DB, uid uint64) (*Blacklist, error) {
	err := db.Debug().Model(Blacklist{}).Where("id = ?", uid).Take(&bl).Error
	if err != nil {
		return &Blacklist{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Blacklist{}, errors.New("Blacklist Not Found")
	}
	if bl.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", &bl.ParentID).Take(&bl.Parent).Error
		if err != nil {
			return &Blacklist{}, err
		}
	}
	return bl, err
}

func (bl *Blacklist) FindBlacklistByParentID(db *gorm.DB, pid uint32) (*[]Blacklist, error) {
	list := []Blacklist{}
	err := db.Debug().Model(&Blacklist{}).Limit(100).Where("parent_id = ?", pid).Find(&list).Error
	if err != nil {
		return &[]Blacklist{}, err
	}
	if len(list) > 0 {
		for i := range list {
			err := db.Debug().Model(&Parent{}).Where("id = ?", list[i].ParentID).Take(&list[i].Parent).Error
			if err != nil {
				return &[]Blacklist{}, err
			}
		}
	}
	return &list, nil
}

func (bl *Blacklist) FindRecordByUrl(db *gorm.DB, link string) (*Blacklist, error) {
	err := db.Debug().Model(&Blacklist{}).Where("url = ?", link).Take(&bl).Error
	if err != nil {
		return &Blacklist{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Blacklist{}, errors.New("Blacklist Not Found")
	}
	return bl, err
}

func (bl *Blacklist) DeleteBlacklist(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Blacklist{}).Where("id = ? and parent_id = ?", pid, uid).Take(&Blacklist{}).Delete(&Blacklist{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Blacklist not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
