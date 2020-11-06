package models

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/repoerna/hms_app/api/utils/hash"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

// Patient ...
type Patient struct {
	ID        uint32    `gorm:"auto_increment" json:"id"`
	SSN       int       `gorm:"primary_key;not null;unique" json:"ssn"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeSave ...
func (p *Patient) BeforeSave() error {
	hashedPassword, err := hash.Hash(p.Password)
	if err != nil {
		return err
	}
	p.Password = string(hashedPassword)
	return nil
}

// Validate ...
func (p *Patient) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if p.Name == "" {
			return errors.New("Required Name")
		}
		if p.SSN == 0 {
			return errors.New("Required SSN")
		}
		if p.Password == "" {
			return errors.New("Required Password")
		}
		if p.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "login":
		if p.Password == "" {
			return errors.New("Required Password")
		}
		if p.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if p.Name == "" {
			return errors.New("Required Name")
		}
		if p.SSN == 0 {
			return errors.New("Required SSN")
		}
		if p.Password == "" {
			return errors.New("Required Password")
		}
		if p.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(p.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

// SavePatient ...
func (p *Patient) SavePatient(db *gorm.DB) (*Patient, error) {

	var err error
	err = db.Debug().Create(&p).Error
	if err != nil {
		return &Patient{}, err
	}
	return p, nil
}

// FindAllPatients ...
func (p *Patient) FindAllPatients(db *gorm.DB) (*[]Patient, error) {
	var err error
	patients := []Patient{}
	err = db.Debug().Model(&Patient{}).Limit(100).Find(&patients).Error
	if err != nil {
		return &[]Patient{}, err
	}
	return &patients, err
}

// FindPatientBySSN ...
func (p *Patient) FindPatientBySSN(db *gorm.DB, ssn uint32) (*Patient, error) {
	var err error
	err = db.Debug().Model(Patient{}).Where("ssn = ?", ssn).Take(&p).Error
	if err != nil {
		return &Patient{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Patient{}, errors.New("User Not Found")
	}
	return p, err
}

// UpdatePatient ...
func (p *Patient) UpdatePatient(db *gorm.DB, ssn uint32) (*Patient, error) {

	// To hash the password
	err := p.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Patient{}).Where("ssn = ?", ssn).Take(&Patient{}).UpdateColumns(
		map[string]interface{}{
			"password":   p.Password,
			"Name":       p.Name,
			"SSN":        p.SSN,
			"email":      p.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Patient{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Patient{}).Where("ssn = ?", ssn).Take(&p).Error
	if err != nil {
		return &Patient{}, err
	}
	return p, nil
}

// DeletePatient ...
func (p *Patient) DeletePatient(db *gorm.DB, ssn uint32) (int64, error) {

	db = db.Debug().Model(&Patient{}).Where("ssn = ?", ssn).Take(&Patient{}).Delete(&Patient{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
