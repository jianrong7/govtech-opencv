package controller

import (
	"govtech-opencv/db"
	"govtech-opencv/model"
	"govtech-opencv/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RegisterReq struct {
	Students []string `json:"students"`
	Teacher  string   `json:"teacher"`
}

// add error handling and test cases
// handle invalid teacher or student use case
func Register(c *fiber.Ctx) error {
	s := new(RegisterReq)
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

type SuspendStudentReq struct {
	Student string `json:"student"`
}

func SuspendStudent(c *fiber.Ctx) error {
	s := new(SuspendStudentReq)
	if err := c.BodyParser(s); err != nil {
		return err
	}
	db.DB.Where("email = ?", s.Student).First(&model.Student{}).Update("is_suspended", true)
	return c.SendStatus(fiber.StatusNoContent)
}

type RetrieveForNotificationsReq struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

func RetrieveForNotifications(c *fiber.Ctx) error {
	var registeredStudentEmails []string
	var validStudentsFromNotification []string
	s := new(RetrieveForNotificationsReq)
	if err := c.BodyParser(s); err != nil {
		return err
	}

	emailsFromNotification := utils.ExtractEmailsFromNotification(s.Notification)
	db.DB.Raw(`SELECT DISTINCT students.email as studentEmail
		FROM students
		WHERE students.is_suspended = FALSE AND students.email IN ?`, emailsFromNotification).Scan(&validStudentsFromNotification)

	db.DB.Raw(`SELECT DISTINCT students.email as studentEmail
		FROM students
		JOIN teacher_students ON students.id = teacher_students.student_id
		JOIN teachers ON teacher_students.teacher_id = teachers.id
		WHERE students.is_suspended = FALSE AND (teachers.email = ? OR students.email IN ?)`, s.Teacher, emailsFromNotification).Scan(&registeredStudentEmails)

	registeredStudentEmails = append(registeredStudentEmails, (validStudentsFromNotification)...)

	return c.JSON(fiber.Map{"recipients": registeredStudentEmails})
}
