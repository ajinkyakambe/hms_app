package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Examination ...
type Examination struct {
	ExaminationID uint32    `gorm:"primary_key;auto_increment" json:"examination_id"`
	AppointmentID uint32    `gorm:"not null" json:"apointment_id"`
	ScheduleCode  string    `gorm:"not null" json:"schedule_code"`
	SSN           int       `gorm:"not null" json:"ssn"`
	EmployeeID    int       `gorm:"not null" json:"employee_id"`
	Anamnesis     string    `gorm:"not null" json:"anamnesis"`
	Diagnosis     string    `gorm:"not null" json:"diagnosis"`
	Prescription  string    `gorm:"not null" json:"prescription"`
	Status        string    `gorm:"default:'available'" json:"status"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Validate ...
func (e *Examination) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if e.AppointmentID == 0 {
			return errors.New("Required Appointment ID")
		}
		if e.ScheduleCode == "" {
			return errors.New("Required Schedule")
		}
		if e.SSN == 0 {
			return errors.New("Required SSN")
		}
		if e.EmployeeID == 0 {
			return errors.New("Required SSN")
		}
		return nil

	default:
		if e.AppointmentID == 0 {
			return errors.New("Required Appointment ID")
		}
		if e.ScheduleCode == "" {
			return errors.New("Required Schedule")
		}
		if e.SSN == 0 {
			return errors.New("Required SSN")
		}
		if e.EmployeeID == 0 {
			return errors.New("Required SSN")
		}
		return nil
	}
}

// SaveExamination ...
func (e *Examination) SaveExamination(db *gorm.DB) (*Examination, error) {

	var err error
	err = db.Debug().Create(&e).Error
	if err != nil {
		return &Examination{}, err
	}
	return e, nil
}

// FindAllExamination ...
func (e *Examination) FindAllExamination(db *gorm.DB) (*[]Examination, error) {
	var err error
	examinations := []Examination{}
	err = db.Debug().Model(&Examination{}).Limit(100).Find(&examinations).Error
	if err != nil {
		return &[]Examination{}, err
	}
	return &examinations, err
}

// FindExaminationByID ...
func (e *Examination) FindExaminationByID(db *gorm.DB, eid uint32) (*Examination, error) {
	var err error
	err = db.Debug().Model(Examination{}).Where("examination_id = ?", eid).Take(&e).Error
	if err != nil {
		return &Examination{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Examination{}, errors.New("User Not Found")
	}
	return e, err
}

// UpdateExamination ...
func (e *Examination) UpdateExamination(db *gorm.DB, eid uint32) (*Examination, error) {

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Debug().Model(&Examination{}).Where("examination_id = ?", eid).Take(&Examination{}).UpdateColumns(
		map[string]interface{}{
			"examination_id": e.ExaminationID,
			"appointment_id": e.AppointmentID,
			"schedule_code":  e.ScheduleCode,
			"ssn":            e.SSN,
			"employee_id":    e.EmployeeID,
			"status":         e.Status,
			"updated_at":     time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return &Examination{}, err
	}

	// This is the display the updated appointment
	if err := tx.Debug().Model(&Examination{}).Where("examination_id = ?", eid).Take(&e).Error; err != nil {
		tx.Rollback()
		return &Examination{}, err
	}

	// Update status on schedule
	getExamination, err := e.FindExaminationByID(tx, eid)
	if err != nil {
		tx.Rollback()
		return &Examination{}, tx.Error
	}

	sch := Schedule{}
	schedule, err := sch.FindSchedulesByCode(tx, getExamination.ScheduleCode)
	if err != nil {
		tx.Rollback()
		return &Examination{}, tx.Error
	}
	schedule.Status = e.Status
	_, err = schedule.UpdateSchedule(tx, getExamination.ScheduleCode)
	if err != nil {
		tx.Rollback()
		return &Examination{}, tx.Error
	}

	tx.Commit()
	return e, nil
}

// DeleteAppointment ...
func (e *Examination) DeleteExamination(db *gorm.DB, aid uint32) (int64, error) {

	db = db.Debug().Model(&Examination{}).Where("appointment = ?", aid).Take(&Examination{}).Delete(&Examination{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
