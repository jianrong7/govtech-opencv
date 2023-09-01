package db

import (
	"fmt"
	"govtech-opencv/app/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var teacherEmails []string = []string{"teacherken@gmail.com", "teacherjoe@gmail.com"}
var studentEmails []string = []string{"studentjon@gmail.com", "studenthon@gmail.com", "studentagnes@gmail.com", "commonstudent1@gmail.com", "commonstudent2@gmail.com", "student_only_under_teacher_ken@gmail.com"}

func seedTeachersTable() {
	for _, email := range teacherEmails {
		teacher := model.Teacher{Email: email}
		DB.Create(&teacher)
	}
}

func seedStudentsTable() {
	for _, email := range studentEmails {
		student := model.Student{Email: email, IsSuspended: false}
		DB.Create(&student)
	}
}

func ConnectToDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	var err error
	dsn := os.Getenv("DB_URI")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatal("Failed to connect to database")
	} else {
		fmt.Println("Successfully connected to the database")
	}

	// remove after dev is done
	DB.Migrator().DropTable(&model.Student{}, &model.Teacher{})
	DB.Migrator().DropTable("teacher_students")

	DB.AutoMigrate(&model.Student{}, &model.Teacher{})
	fmt.Println("Database migrated")

	// remove after dev is done
	seedTeachersTable()
	seedStudentsTable()
}
