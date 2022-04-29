package models

import (
	"errors"
	"html"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Parent struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nama      string    `gorm:"size:255;not null;" json:"nama"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	FCM       string    `gorm:"type:text" json:"fcm"`
	IsChange  bool      `gorm:"default:false" json:"is_change"`
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
	p.IsChange = false
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

	case "reset":
		if p.Email == "" {
			return errors.New("butuh email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "newpassword":
		if p.Password == "" {
			return errors.New("butuh password")
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
			"is_change":  false,
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

func (p *Parent) UpdateParentFCM(db *gorm.DB, uid uint32) (*Parent, error) {
	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).UpdateColumns(
		map[string]interface{}{
			"fcm":        p.FCM,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Parent{}, db.Error
	}
	// This is the display the updated parent
	err := db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&p).Error
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

func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var lowerCharSet = "abcdedfghijklmnopqrst"
	var upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var specialCharSet = "!@#$%&*"
	var numberSet = "0123456789"
	var allCharSet = lowerCharSet + upperCharSet + specialCharSet + numberSet
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func (p *Parent) ResetParentPassword(db *gorm.DB, uid uint32) (*Parent, string, error) {
	rand.Seed(time.Now().Unix())
	minSpecialChar := 1
	minNum := 1
	minUpperCase := 1
	passwordLength := 8
	password := GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase)

	p.Password = password
	err := p.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&Parent{}).UpdateColumns(
		map[string]interface{}{
			"password":   p.Password,
			"is_change":  true,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Parent{}, password, db.Error
	}

	// This is the display the updated parent
	err = db.Debug().Model(&Parent{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Parent{}, password, err
	}
	return p, password, nil
}
