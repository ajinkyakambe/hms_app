package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/repoerna/hms_app/api/models"
)

var patients = []models.Patient{
	models.Patient{
		Name:     "Purna",
		SSN:      1234567777,
		Email:    "purna@gmail.com",
		Password: "password",
	},
	models.Patient{
		Name:     "Fajar",
		SSN:      1234568888,
		Email:    "Fajar@gmail.com",
		Password: "password",
	},
}

var employees = []models.Employee{
	models.Employee{
		Name:       "dr. Bob",
		EmployeeID: 201103001,
		Email:      "bob@gmail.com",
		Password:   "password",
		Department: "Dokter Umum",
	},
	models.Employee{
		Name:       "Ella",
		EmployeeID: 201103002,
		Email:      "ella@gmail.com",
		Password:   "password",
		Department: "Perawat",
	},
}

var schedules = []models.Schedule{
	models.Schedule{
		ScheduleCode: "SC-00001",
		Day:          "Senin",
		StartTime:    "08:00",
		EndTime:      "11:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00002",
		Day:          "Senin",
		StartTime:    "14:00",
		EndTime:      "16:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00003",
		Day:          "Senin",
		StartTime:    "20:00",
		EndTime:      "21:30",
	},
	models.Schedule{
		ScheduleCode: "SC-00004",
		Day:          "Selasa",
		StartTime:    "10:00",
		EndTime:      "12:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00005",
		Day:          "Selasa",
		StartTime:    "15:00",
		EndTime:      "18:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00006",
		Day:          "Rabu",
		StartTime:    "13:00",
		EndTime:      "18:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00007",
		Day:          "Kamis",
		StartTime:    "08:00",
		EndTime:      "11:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00008",
		Day:          "Jumat",
		StartTime:    "14:00",
		EndTime:      "18:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00009",
		Day:          "Sabtu",
		StartTime:    "08:00",
		EndTime:      "11:00",
	},
	models.Schedule{
		ScheduleCode: "SC-00010",
		Day:          "Minggu",
		StartTime:    "20:00",
		EndTime:      "22:00",
	},
}

// Load ...
func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Patient{}, &models.Employee{}, &models.Schedule{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Patient{}, &models.Employee{}, &models.Schedule{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range patients {
		err = db.Debug().Model(&models.Patient{}).Create(&patients[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	for i := range employees {
		err = db.Debug().Model(&models.Employee{}).Create(&employees[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	for i := range schedules {
		err = db.Debug().Model(&models.Schedule{}).Create(&schedules[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

}
