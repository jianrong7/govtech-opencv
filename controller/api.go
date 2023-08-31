package controller

import (
	"govtech-opencv/db"
	"govtech-opencv/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RegisterStruct struct {
	Students []string `json:"students"`
	Teacher  string   `json:"teacher"`
}

// add error handling and test cases
// handle invalid teacher or student use case
func Register(c *fiber.Ctx) error {
	s := new(RegisterStruct)
	if err := c.BodyParser(s); err != nil {
		return err
	}
	teacherEmail := s.Teacher
	var teacher model.Teacher
	db.DB.Where("email = ?", teacherEmail).First(&teacher)

	for _, studentEmail := range s.Students {
		var student model.Student
		db.DB.Where("email = ?", studentEmail).First(&student)
		student.Teachers = append(student.Teachers, teacher)

		db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&student)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetCommonStudents(c *fiber.Ctx) error {
	var commonStudentEmails []string

	queryString := string(c.Request().URI().QueryString())
	queryString = strings.Replace(queryString, "teacher=", "", -1)
	queryString = strings.Replace(queryString, "%40", "@", -1)
	teacherEmails := strings.Split(queryString, "&")

	db.DB.Raw(`SELECT DISTINCT students.email as studentEmail
		FROM students
		JOIN teacher_students ON students.id = teacher_students.student_id
		JOIN teachers ON teacher_students.teacher_id = teachers.id
		WHERE teachers.email in ?`, teacherEmails).Scan(&commonStudentEmails)

	return c.JSON(fiber.Map{"students": commonStudentEmails})
}

type SuspendStudentStruct struct {
	Student string `json:"student"`
}

func SuspendStudent(c *fiber.Ctx) error {
	s := new(SuspendStudentStruct)
	if err := c.BodyParser(s); err != nil {
		return err
	}
	db.DB.Where("email = ?", s.Student).First(&model.Student{}).Update("isSuspended", true)
	return c.SendStatus(fiber.StatusNoContent)
}

func RetrieveForNotifications(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "RetrieveForNotifications", "data": nil})
}
