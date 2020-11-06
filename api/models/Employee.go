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

// Employee ...
type Employee struct {
	EmployeeID int       `gorm:"primary_key;auto_increment" json:"employee_id"`
	Name       string    `gorm:"size:255;not null;unique" json:"name"`
	Email      string    `gorm:"size:100;not null;unique" json:"email"`
	Password   string    `gorm:"size:100;not null;" json:"password"`
	Department string    `gorm:"size:100;not null;" json:"department"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeSave ...
func (e *Employee) BeforeSave() error {
	hashedPassword, err := hash.Hash(e.Password)
	if err != nil {
		return err
	}
	e.Password = string(hashedPassword)
	return nil
}

// Validate ...
func (e *Employee) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if e.Name == "" {
			return errors.New("Required Name")
		}
		if e.Department == "" {
			return errors.New("Required Department")
		}
		if e.Password == "" {
			return errors.New("Required Password")
		}
		if e.EmployeeID == 0 {
			return errors.New("Required Employee ID")
		}
		if e.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(e.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "login":
		if e.Password == "" {
			return errors.New("Required Password")
		}
		if e.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(e.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if e.Name == "" {
			return errors.New("Required Name")
		}
		if e.Department == "" {
			return errors.New("Required Department")
		}
		if e.Password == "" {
			return errors.New("Required Password")
		}
		if e.EmployeeID == 0 {
			return errors.New("Required Employee ID")
		}
		if e.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(e.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

// SaveEmployee ...
func (e *Employee) SaveEmployee(db *gorm.DB) (*Employee, error) {

	var err error
	err = db.Debug().Create(&e).Error
	if err != nil {
		return &Employee{}, err
	}
	return e, nil
}

// FindAllEmployee ...
func (e *Employee) FindAllEmployee(db *gorm.DB) (*[]Employee, error) {
	var err error
	employees := []Employee{}
	err = db.Debug().Model(&Employee{}).Limit(100).Find(&employees).Error
	if err != nil {
		return &[]Employee{}, err
	}
	return &employees, err
}

// FindEmployeeByID ...
func (e *Employee) FindEmployeeByID(db *gorm.DB, employeeID int) (*Employee, error) {
	var err error
	err = db.Debug().Model(Employee{}).Where("employee_id = ?", employeeID).Take(&e).Error
	if err != nil {
		return &Employee{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Employee{}, errors.New("User Not Found")
	}
	return e, err
}

// UpdateEmployee ...
func (e *Employee) UpdateEmployee(db *gorm.DB, employeeID uint32) (*Employee, error) {

	// To hash the password
	err := e.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Employee{}).Where("employee_id=?", employeeID).Take(&Employee{}).UpdateColumns(
		map[string]interface{}{
			"password":    e.Password,
			"name":        e.Name,
			"employee_id": e.EmployeeID,
			"department":  e.Department,
			"email":       e.Email,
			"updated_at":  time.Now(),
		},
	)
	if db.Error != nil {
		return &Employee{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Employee{}).Where("employee_id = ?", employeeID).Take(&e).Error
	if err != nil {
		return &Employee{}, err
	}
	return e, nil
}

// DeleteEmployee ...
func (e *Employee) DeleteEmployee(db *gorm.DB, employeeID uint32) (int64, error) {

	db = db.Debug().Model(&Employee{}).Where("employee_id = ?", employeeID).Take(&Employee{}).Delete(&Employee{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
