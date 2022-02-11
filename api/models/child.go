package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Child struct {
	ID        	uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Nama     	string    	`gorm:"size:255;not null;unique" json:"nama"`
	Email   	string    	`gorm:"size:255;not null;" json:"email"`
	Password   	string    	`gorm:"size:255;not null;" json:"password"`
	Parent    	Parent     `json:"Parent"`
	ParentID  	uint32    `gorm:"not null" json:"parent_id"`
	CreatedAt 	time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func HashChild(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPasswordChild(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (c *Child) BeforeSaveChild() error {
	hashedPassword, err := Hash(c.Password)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	return nil
}

func (c *Child) Prepare() {
	c.ID = 0
	c.Nama = html.EscapeString(strings.TrimSpace(c.Nama))
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.Parent = Parent{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Child) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.Nama == "" {
			return errors.New("required nama")
		}
		if c.Password == "" {
			return errors.New("required password")
		}
		if c.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	case "login":
		if c.Password == "" {
			return errors.New("required password")
		}
		if c.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	default:
		if c.Nama == "" {
			return errors.New("required nama")
		}
		if c.Password == "" {
			return errors.New("required password")
		}
		if c.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}


func (c *Child) SaveChild(db *gorm.DB) (*Child, error) {
	var err error
	err = db.Debug().Model(&Child{}).Create(&c).Error
	if err != nil {
		return &Child{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", c.ParentID).Take(&c.Parent).Error
		if err != nil {
			return &Child{}, err
		}
	}
	return c, nil
}

func (c *Child) FindAllChilds(db *gorm.DB) (*[]Child, error) {
	var err error
	childs := []Child{}
	err = db.Debug().Model(&Child{}).Limit(100).Find(&childs).Error
	if err != nil {
		return &[]Child{}, err
	}
	if len(childs) > 0 {
		for i, _ := range childs {
			err := db.Debug().Model(&Parent{}).Where("id = ?", childs[i].ParentID).Take(&childs[i].Parent).Error
			if err != nil {
				return &[]Child{}, err
			}
		}
	}
	return &childs, nil
}

func (c *Child) FindChildByID(db *gorm.DB, pid uint64) (*Child, error) {
	var err error
	err = db.Debug().Model(&Child{}).Where("id = ?", pid).Take(&c).Error
	if err != nil {
		return &Child{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", c.ParentID).Take(&c.Parent).Error
		if err != nil {
			return &Child{}, err
		}
	}
	return c, nil
}

func (c *Child) UpdateAChild(db *gorm.DB) (*Child, error) {

	var err error

	err = db.Debug().Model(&Child{}).Where("id = ?", c.ID).Updates(Child{Nama: c.Nama, Email: c.Email, Password: c.Password, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Child{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&Parent{}).Where("id = ?", c.ParentID).Take(&c.Parent).Error
		if err != nil {
			return &Child{}, err
		}
	}
	return c, nil
}

func (c *Child) DeleteAChild(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Child{}).Where("id = ? and parent_id = ?", pid, uid).Take(&Child{}).Delete(&Child{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Child not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}