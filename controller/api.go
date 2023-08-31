package controller

import (
	"govtech-opencv/db"
	"govtech-opencv/model"
	"govtech-opencv/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// add test cases
func Register(c *fiber.Ctx) error {
	type RegisterReq struct {
		Students []string `json:"students"`
		Teacher  string   `json:"teacher"`
	}

	s := new(RegisterReq)
	if err := c.BodyParser(s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	teacherEmail := s.Teacher
	var teacher model.Teacher
	err := db.DB.Where("email = ?", teacherEmail).First(&teacher).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "teacher not found"})
	}

	for _, studentEmail := range s.Students {
		var student model.Student
		err = db.DB.Where("email = ?", studentEmail).First(&student).Error

		if err != nil {
			student = model.Student{
				Email: studentEmail,
				Teachers: []model.Teacher{
					teacher,
				},
			}
			db.DB.Create(&student)
		} else {
			student.Teachers = append(student.Teachers, teacher)
			db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&student)
		}

	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetCommonStudents(c *fiber.Ctx) error {
	var commonStudentEmails []string

	queryString := string(c.Request().URI().QueryString())
	queryString = strings.Replace(queryString, "teacher=", "", -1)
	queryString = strings.Replace(queryString, "%40", "@", -1)
	teacherEmails := strings.Split(queryString, "&")

	if len(teacherEmails) == 0 {
		return c.JSON(fiber.Map{"students": commonStudentEmails})
	}

	db.DB.Raw(`SELECT DISTINCT students.email as studentEmail
		FROM students
		JOIN teacher_students ON students.id = teacher_students.student_id
		JOIN teachers ON teacher_students.teacher_id = teachers.id
		WHERE teachers.email in ?
		AND students.deleted_at IS NULL
		AND teachers.deleted_at IS NULL`, teacherEmails).Scan(&commonStudentEmails)

	return c.JSON(fiber.Map{"students": commonStudentEmails})
}

func SuspendStudent(c *fiber.Ctx) error {
	type SuspendStudentReq struct {
		Student string `json:"student"`
	}

	s := new(SuspendStudentReq)
	if err := c.BodyParser(s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	err := db.DB.Where("email = ?", s.Student).First(&model.Student{}).Update("is_suspended", true).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "student not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func RetrieveForNotifications(c *fiber.Ctx) error {
	type RetrieveForNotificationsReq struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}

	var registeredStudentEmails []string
	var validStudentsFromNotification []string

	s := new(RetrieveForNotificationsReq)
	if err := c.BodyParser(s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	emailsFromNotification := utils.ExtractEmailsFromNotification(s.Notification)
	if len(emailsFromNotification) > 0 {
		db.DB.Raw(`SELECT DISTINCT students.email as studentEmail
		FROM students
		WHERE students.is_suspended = FALSE
		AND students.email IN ?
		AND students.deleted_at IS NULL`, emailsFromNotification).Scan(&validStudentsFromNotification)
	}

	db.DB.Raw(`SELECT DISTINCT students.email as studentEmail
		FROM students
		JOIN teacher_students ON students.id = teacher_students.student_id
		JOIN teachers ON teacher_students.teacher_id = teachers.id
		WHERE students.is_suspended = FALSE
		AND (teachers.email = ? OR students.email IN ?)
		AND students.deleted_at IS NULL
		AND teachers.deleted_at IS NULL`, s.Teacher, emailsFromNotification).Scan(&registeredStudentEmails)

	registeredStudentEmails = append(registeredStudentEmails, (validStudentsFromNotification)...)

	return c.JSON(fiber.Map{"recipients": registeredStudentEmails})
}
