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
	Nama      string    `gorm:"size:255;not null;" json:"nama"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	LastLogin time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_login"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (p *Parent) BeforeSave() error {
	hashedPassword, err := Hash(p.Password)
	if err != nil {
		return err
	}
	p.Password = string(hashedPassword)
	return nil
}

func (p *Parent) Prepare() {
	p.ID = 0
	p.Nama = html.EscapeString(strings.TrimSpace(p.Nama))
	p.Email = html.EscapeString(strings.TrimSpace(p.Email))
	p.LastLogin = time.Now()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Parent) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if p.Nama == "" {
			return errors.New("butuh nama")
		}
		if p.Password == "" {
			return errors.New("butuh password")
		}
		if p.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "login":
		if p.Password == "" {
			return errors.New("butuh password")
		}
		if p.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "updatepassword":
		if p.Password == "" {
			return errors.New("butuh password")
		}
		if p.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "updateprofile":
		if p.Nama == "" {
			return errors.New("butuh nama")
		}
		if p.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	default:
		if p.Nama == "" {
			return errors.New("butuh nama")
		}
		if p.Password == "" {
			return errors.New("butuh password")
		}
		if p.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

func (p *Parent) SaveParent(db *gorm.DB) (*Parent, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Parent{}, err
	}
	return p, nil
}

func (p *Parent) FindAllParents(db *gorm.DB) (*[]Parent, error) {
	parents := []Parent{}
	err := db.Debug().Model(&Parent{}).Limit(100).Find(&parents).Error
	if err != nil {
		return &[]Parent{}, err
	}
	return &parents, err
}

func (p *Parent) FindParentByID(db *gorm.DB, uid uint32) (*Parent, error) {
	err := db.Debug().Model(Parent{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Parent{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Parent{}, errors.New("Parent Not Found")
	}
	return p, err
}

func (p *Parent) UpdateAParent(db *gorm.DB, uid uint32) (*Parent, error) {
	// To hash the password
	err := p.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).UpdateColumns(
		map[string]interface{}{
			"password":   p.Password,
			"nama":       p.Nama,
			"email":      p.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Parent{}, db.Error
	}
	// This is the display the updated parent
	err = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Parent{}, err
	}
	return p, nil
}

func (p *Parent) UpdateParentProfile(db *gorm.DB, uid uint32) (*Parent, error) {
	err := p.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).UpdateColumns(
		map[string]interface{}{
			"nama":       p.Nama,
			"email":      p.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Parent{}, db.Error
	}
	// This is the display the updated parent
	err = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Parent{}, err
	}
	return p, nil
}

func (p *Parent) UpdateParentPassword(db *gorm.DB, uid uint32) (*Parent, error) {
	err := p.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).UpdateColumns(
		map[string]interface{}{
			"password":   p.Password,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Parent{}, db.Error
	}
	// This is the display the updated parent
	err = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Parent{}, err
	}
	return p, nil
}

func (p *Parent) DeleteAParent(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).Delete(&Parent{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
