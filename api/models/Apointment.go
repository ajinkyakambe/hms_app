package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Appointment ...
type Appointment struct {
	AppointmentID uint32    `gorm:"primary_key;auto_increment" json:"apointment_id"`
	ScheduleCode  string    `gorm:"not null" json:"schedule_code"`
	SSN           int       `gorm:"not null" json:"ssn"`
	EmployeeID    int       `gorm:"not null" json:"employee_id"`
	StartTime     string    `gorm:"size:100;not null;" json:"start_time"`
	EndTime       string    `gorm:"size:100;not null;" json:"end_time"`
	Status        string    `gorm:"default:'available'" json:"status"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Validate ...
func (a *Appointment) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if a.AppointmentID == 0 {
			return errors.New("Required Appointment ID")
		}
		if a.ScheduleCode == "" {
			return errors.New("Required Schedule")
		}
		if a.StartTime == "" {
			return errors.New("Required Start Time")
		}
		if a.EndTime == "" {
			return errors.New("Required End Time")
		}
		if a.SSN == 0 {
			return errors.New("Required SSN")
		}
		return nil

	default:
		if a.AppointmentID == 0 {
			return errors.New("Required Appointment ID")
		}
		if a.ScheduleCode == "" {
			return errors.New("Required Schedule")
		}
		if a.StartTime == "" {
			return errors.New("Required Start Time")
		}
		if a.EndTime == "" {
			return errors.New("Required End Time")
		}
		if a.SSN == 0 {
			return errors.New("Required SSN")
		}
		return nil
	}
}

// SaveAppointment ...
func (a *Appointment) SaveAppointment(db *gorm.DB) (*Appointment, error) {

	var err error
	err = db.Debug().Create(&a).Error
	if err != nil {
		return &Appointment{}, err
	}
	return a, nil
}

// FindAllAppointment ...
func (a *Appointment) FindAllAppointment(db *gorm.DB) (*[]Appointment, error) {
	var err error
	appointments := []Appointment{}
	err = db.Debug().Model(&Appointment{}).Limit(100).Find(&appointments).Error
	if err != nil {
		return &[]Appointment{}, err
	}
	return &appointments, err
}

// FindAppointmentByID ...
func (a *Appointment) FindAppointmentByID(db *gorm.DB, aid uint32) (*Appointment, error) {
	var err error
	err = db.Debug().Model(Appointment{}).Where("appointment_id = ?", aid).Take(&a).Error
	if err != nil {
		return &Appointment{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Appointment{}, errors.New("User Not Found")
	}
	return a, err
}

// UpdateAppointment ...
func (a *Appointment) UpdateAppointment(db *gorm.DB, aid uint32) (*Appointment, error) {

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Debug().Model(&Appointment{}).Where("appointment_id = ?", aid).Take(&Appointment{}).UpdateColumns(
		map[string]interface{}{
			"appointment_id": a.AppointmentID,
			"schedule_code":  a.ScheduleCode,
			"ssn":            a.SSN,
			"employee_id":    a.EmployeeID,
			"status":         a.Status,
			"updated_at":     time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return &Appointment{}, err
	}

	// This is the display the updated appointment
	if err := tx.Debug().Model(&Appointment{}).Where("appointment_id = ?", aid).Take(&a).Error; err != nil {
		tx.Rollback()
		return &Appointment{}, err
	}

	// Update status on schedule
	getAppointment, err := a.FindAppointmentByID(tx, aid)
	if err != nil {
		tx.Rollback()
		return &Appointment{}, tx.Error
	}

	sch := Schedule{}
	schedule, err := sch.FindSchedulesByCode(tx, getAppointment.ScheduleCode)
	if err != nil {
		tx.Rollback()
		return &Appointment{}, tx.Error
	}
	schedule.Status = a.Status
	_, err = schedule.UpdateSchedule(tx, getAppointment.ScheduleCode)
	if err != nil {
		tx.Rollback()
		return &Appointment{}, tx.Error
	}

	tx.Commit()
	return a, nil
}

// DeleteAppointment ...
func (a *Appointment) DeleteAppointment(db *gorm.DB, aid uint32) (int64, error) {

	db = db.Debug().Model(&Appointment{}).Where("appointment_id = ?", aid).Take(&Appointment{}).Delete(&Appointment{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
