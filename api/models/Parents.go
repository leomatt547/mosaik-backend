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

type Parent struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nama      string    `gorm:"size:255;not null;unique" json:"nama"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *Parent) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *Parent) Prepare() {
	u.ID = 0
	u.Nama = html.EscapeString(strings.TrimSpace(u.Nama))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (p *Parent) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if p.Nama == "" {
			return errors.New("required nama")
		}
		if p.Password == "" {
			return errors.New("required password")
		}
		if p.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}

		return nil
	case "login":
		if p.Password == "" {
			return errors.New("required Password")
		}
		if p.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	default:
		if p.Nama == "" {
			return errors.New("required nama")
		}
		if p.Password == "" {
			return errors.New("required password")
		}
		if p.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

func (u *Parent) SaveParent(db *gorm.DB) (*Parent, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Parent{}, err
	}
	return u, nil
}

func (u *Parent) FindAllParents(db *gorm.DB) (*[]Parent, error) {
	var err error
	parents := []Parent{}
	err = db.Debug().Model(&Parent{}).Limit(100).Find(&parents).Error
	if err != nil {
		return &[]Parent{}, err
	}
	return &parents, err
}

func (u *Parent) FindParentByID(db *gorm.DB, uid uint32) (*Parent, error) {
	var err error
	err = db.Debug().Model(Parent{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Parent{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Parent{}, errors.New("Parent Not Found")
	}
	return u, err
}

func (u *Parent) UpdateAParent(db *gorm.DB, uid uint32) (*Parent, error) {
	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nama":  u.Nama,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Parent{}, db.Error
	}
	// This is the display the updated parent
	err = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Parent{}, err
	}
	return u, nil
}

func (u *Parent) DeleteAParent(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).Delete(&Parent{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}