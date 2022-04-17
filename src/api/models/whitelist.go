package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Whitelist struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Url       string    `gorm:"uniqueIndex:whitelists_url_parent_id;type:text;not null;" json:"url"`
	ParentID  uint32    `gorm:"uniqueIndex:whitelists_url_parent_id;not null" json:"parent_id"`
	Parent    Parent    `json:"Parent"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (wl *Whitelist) Validate() error {
	if wl.Url == "" {
		return errors.New("butuh url")
	}
	return nil
}

func (wl *Whitelist) SaveWhitelist(db *gorm.DB) (*Whitelist, error) {
	err := db.Debug().Create(&wl).Error
	if err != nil {
		return &Whitelist{}, err
	}
	if wl.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", &wl.ParentID).Take(&wl.Parent).Error
		if err != nil {
			return &Whitelist{}, err
		}
	}
	return wl, err
}

func (wl *Whitelist) FindAllWhitelist(db *gorm.DB) (*[]Whitelist, error) {
	list := []Whitelist{}
	err := db.Debug().Model(&Whitelist{}).Limit(100).Find(&list).Error
	if err != nil {
		return &[]Whitelist{}, err
	}
	if len(list) > 0 {
		for i := range list {
			err := db.Debug().Model(&Parent{}).Where("id = ?", list[i].ParentID).Take(&list[i].Parent).Error
			if err != nil {
				return &[]Whitelist{}, err
			}
		}
	}
	return &list, err
}

func (wl *Whitelist) FindWhitelistByID(db *gorm.DB, uid uint64) (*Whitelist, error) {
	err := db.Debug().Model(Whitelist{}).Where("id = ?", uid).Take(&wl).Error
	if err != nil {
		return &Whitelist{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Whitelist{}, errors.New("Whitelist Not Found")
	}
	if wl.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", &wl.ParentID).Take(&wl.Parent).Error
		if err != nil {
			return &Whitelist{}, err
		}
	}
	return wl, err
}

func (wl *Whitelist) FindWhitelistByParentID(db *gorm.DB, pid uint32) (*[]Whitelist, error) {
	list := []Whitelist{}
	err := db.Debug().Model(&Whitelist{}).Limit(100).Where("parent_id = ?", pid).Find(&list).Error
	if err != nil {
		return &[]Whitelist{}, err
	}
	if len(list) > 0 {
		for i := range list {
			err := db.Debug().Model(&Parent{}).Where("id = ?", list[i].ParentID).Take(&list[i].Parent).Error
			if err != nil {
				return &[]Whitelist{}, err
			}
		}
	}
	return &list, nil
}

func (wl *Whitelist) FindRecordByUrl(db *gorm.DB, link string, cid uint64) (*Whitelist, error) {
	child := Child{}
	err := db.Debug().Model(&Child{}).Where("id = ?", cid).Take(&child).Error
	if err != nil {
		return &Whitelist{}, err
	}
	err = db.Debug().Model(&Whitelist{}).Where("url = ? and parent_id = ?", link, child.ParentID).Take(&wl).Error
	if err != nil {
		return &Whitelist{}, err
	}
	err = db.Debug().Model(&Parent{}).Where("id = ?", wl.ParentID).Take(&wl.Parent).Error
	if err != nil {
		return &Whitelist{}, err
	}
	return wl, err
}

func (wl *Whitelist) DeleteWhitelist(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Whitelist{}).Where("id = ? and parent_id = ?", pid, uid).Take(&Whitelist{}).Delete(&Whitelist{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Whitelist not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
