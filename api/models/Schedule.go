package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Schedule ...
type Schedule struct {
	ID           uint32    `gorm:"auto_increment" json:"id"`
	ScheduleCode string    `gorm:"primary_key" json:"schedule_code"`
	Day          string    `gorm:"size:100;not null;" json:"day"`
	StartTime    string    `gorm:"size:100;not null;" json:"start_time"`
	EndTime      string    `gorm:"size:100;not null;" json:"end_time"`
	Status       string    `gorm:"default:'available'" json:"status"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Validate ...
func (s *Schedule) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if s.ScheduleCode == "" {
			return errors.New("Required Schedule Code")
		}
		if s.Day == "" {
			return errors.New("Required Day Name")
		}
		if s.StartTime == "" {
			return errors.New("Required Start Time")
		}
		if s.EndTime == "" {
			return errors.New("Required End Time")
		}
		return nil

	default:
		if s.ScheduleCode == "" {
			return errors.New("Required Schedule Code")
		}
		if s.Day == "" {
			return errors.New("Required Day Name")
		}
		if s.StartTime == "" {
			return errors.New("Required Start Time")
		}
		if s.EndTime == "" {
			return errors.New("Required End Time")
		}
		return nil
	}
}

// SaveSchedule ...
func (s *Schedule) SaveSchedule(db *gorm.DB) (*Schedule, error) {
	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return &Schedule{}, err
	}
	return s, nil
}

// FindSchedulesByCode ...
func (s *Schedule) FindSchedulesByCode(db *gorm.DB, sc string) (*Schedule, error) {
	var err error
	err = db.Debug().Model(Schedule{}).Where("sc = ?", sc).Take(&s).Error
	if err != nil {
		return &Schedule{}, err
	}
	return s, err
}

// FindAllSchedules ...
func (s *Schedule) FindAllSchedules(db *gorm.DB) (*[]Schedule, error) {
	var err error
	schedules := []Schedule{}
	err = db.Debug().Model(&Schedule{}).Limit(100).Find(&schedules).Error
	if err != nil {
		return &[]Schedule{}, err
	}
	return &schedules, err
}

// UpdateSchedule ...
func (s *Schedule) UpdateSchedule(db *gorm.DB, sc string) (*Schedule, error) {

	db = db.Debug().Model(&Schedule{}).Where("schedule_code=?", sc).Take(&Schedule{}).UpdateColumns(
		map[string]interface{}{
			"id":            s.ID,
			"schedule_code": s.ScheduleCode,
			"day":           s.Day,
			"status":        s.Status,
			"start_time":    s.StartTime,
			"end_time":      s.EndTime,
			"updated_at":    time.Now(),
		},
	)
	if db.Error != nil {
		return &Schedule{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&Schedule{}).Where("schedule_code = ?", sc).Take(&s).Error
	if err != nil {
		return &Schedule{}, err
	}
	return s, nil
}

// DeleteSchedule ...
func (s *Schedule) DeleteSchedule(db *gorm.DB, sc string) (int64, error) {

	db = db.Debug().Model(&Schedule{}).Where("schedule_code = ?", sc).Take(&Schedule{}).Delete(&Schedule{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
