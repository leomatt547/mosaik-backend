package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Child struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Nama      string    `gorm:"size:255;not null;" json:"nama"`
	Email     string    `gorm:"size:100;not null;unique;" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Parent    Parent    `json:"Parent"`
	ParentID  uint32    `gorm:"not null" json:"parent_id"`
	LastLogin time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_login"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func HashChild(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPasswordChild(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (c *Child) BeforeSave() error {
	hashedPassword, err := HashChild(c.Password)
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
	c.LastLogin = time.Now()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Child) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.Nama == "" {
			return errors.New("butuh nama")
		}
		if c.Password == "" {
			return errors.New("butuh password")
		}
		if c.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "login":
		if c.Password == "" {
			return errors.New("butuh password")
		}
		if c.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "updatepassword":
		if c.Password == "" {
			return errors.New("butuh password")
		}
		if c.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "updateprofile":
		if c.Nama == "" {
			return errors.New("butuh nama")
		}
		if c.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	default:
		if c.Nama == "" {
			return errors.New("butuh nama")
		}
		if c.Password == "" {
			return errors.New("butuh password")
		}
		if c.Email == "" {
			return errors.New("butuh email")
		}
		if c.ParentID == 0 {
			return errors.New("butuh parent_id")
		}
		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

func (c *Child) SaveChild(db *gorm.DB) (*Child, error) {
	var err error
	//err = db.Debug().Model(&Child{}).Create(&c).Error
	err = db.Debug().Create(&c).Error
	if err != nil {
		return &Child{}, err
	}
	if c.ID != 0 {
		//Dapatkan id Parent apakah ada atau tidak
		err = db.Debug().Model(&Parent{}).Where("id = ?", &c.ParentID).Take(&c.Parent).Error
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
		for i := range childs {
			err := db.Debug().Model(&Parent{}).Where("id = ?", childs[i].ParentID).Take(&childs[i].Parent).Error
			if err != nil {
				return &[]Child{}, err
			}
		}
	}
	return &childs, nil
}

func (c *Child) FindChildByID(db *gorm.DB, cid uint64) (*Child, error) {
	var err error
	err = db.Debug().Model(&Child{}).Where("id = ?", cid).Take(&c).Error
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

func (c *Child) FindChildbyParentID(db *gorm.DB, pid uint32) (*[]Child, error) {
	var err error
	child := []Child{}
	err = db.Debug().Model(&Child{}).Limit(100).Where("parent_id = ?", pid).Find(&child).Error
	if err != nil {
		return &[]Child{}, err
	}
	if len(child) > 0 {
		//Dapatkan id Parent
		for i := range child {
			err := db.Debug().Model(&Parent{}).Where("id = ?", pid).Take(&child[i].Parent).Error
			if err != nil {
				return &[]Child{}, err
			}
		}
	}
	return &child, nil
}

func (c *Child) UpdateAChild(db *gorm.DB, uid uint64) (*Child, error) {
	err := c.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Child{}).Where("id = ?", uid).Take(&Child{}).UpdateColumns(
		map[string]interface{}{
			"nama":       c.Nama,
			"email":      c.Email,
			"password":   c.Password,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Child{}, db.Error
	}
	// This is the display the updated child
	err = db.Debug().Model(&Child{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Child{}, err
	}
	return c, nil
}

func (c *Child) UpdateChildProfile(db *gorm.DB, uid uint64) (*Child, error) {
	err := c.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Child{}).Where("id = ?", uid).Take(&Child{}).UpdateColumns(
		map[string]interface{}{
			"nama":       c.Nama,
			"email":      c.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Child{}, db.Error
	}
	// This is the display the updated child
	err = db.Debug().Model(&Child{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Child{}, err
	}
	return c, nil
}

func (c *Child) UpdateChildPassword(db *gorm.DB, uid uint64) (*Child, error) {

	err := c.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Child{}).Where("id = ?", uid).Take(&Child{}).UpdateColumns(
		map[string]interface{}{
			"password":   c.Password,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Child{}, db.Error
	}
	// This is the display the updated child
	err = db.Debug().Model(&Child{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Child{}, err
	}
	err = db.Debug().Model(&Parent{}).Where("id = ?", c.ParentID).Take(&c.Parent).Error
	if err != nil {
		return &Child{}, err
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
