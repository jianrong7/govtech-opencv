package controller

import (
	"govtech-opencv/app/dto"
	"govtech-opencv/app/model"
	"govtech-opencv/db"
	"govtech-opencv/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register students to a teacher. If student does not exist, create it.
// If teacher does not exist, return 404.
// @Description Register students to a teacher.
// @Summary Register students to a teacher.
// @Accept json
// @Produce json
// @Param register body dto.RegisterReq true "Register"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 404 {string} string "teacher not found"
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	s := new(dto.RegisterReq)
	if err := c.BodyParser(s); err != nil || len(s.Students) == 0 || s.Teacher == "" {
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

// Get common students from a list of teachers.
// @Description Get common students from a list of teachers.
// @Summary Get common students from a list of teachers.
// @Produce json
// @Param teacher query string true "Teacher"
// @Success 200 {object} string "students"
// @Failure 400 {string} string "invalid request"
// @Router /api/commonstudents [get]
func GetCommonStudents(c *fiber.Ctx) error {
	var commonStudentEmails []string

	queryString := string(c.Request().URI().QueryString())
	if queryString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "no teacher indicated"})
	}
	queryString = strings.Replace(queryString, "teacher=", "", -1)
	queryString = strings.Replace(queryString, "%40", "@", -1)
	teacherEmails := strings.Split(queryString, "&")

	db.DB.Raw(`SELECT students.email as studentEmail
		FROM students
		JOIN teacher_students ON students.id = teacher_students.student_id
		JOIN teachers ON teacher_students.teacher_id = teachers.id
		WHERE teachers.email IN ?
		AND students.deleted_at IS NULL
		AND teachers.deleted_at IS NULL
		GROUP BY studentEmail
		HAVING COUNT(DISTINCT teachers.email) = ?`, teacherEmails, len(teacherEmails)).Scan(&commonStudentEmails)

	if len(commonStudentEmails) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no students found"})
	} else {
		return c.JSON(fiber.Map{"students": commonStudentEmails})
	}
}

// Suspend a student. If student does not exist, return 404.
// @Description Suspend a student.
// @Summary Suspend a student.
// @Accept json
// @Produce json
// @Param suspend body dto.SuspendStudentReq true "Suspend"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 404 {string} string "student not found"
// @Router /api/suspend [post]
func SuspendStudent(c *fiber.Ctx) error {
	s := new(dto.SuspendStudentReq)
	if err := c.BodyParser(s); err != nil || s.Student == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	err := db.DB.Where("email = ?", s.Student).First(&model.Student{}).Update("is_suspended", true).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "student not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Retrieve students who can receive notifications from a teacher.
// @Description Retrieve students who can receive notifications from a teacher.
// @Summary Retrieve students who can receive notifications from a teacher.
// @Accept json
// @Produce json
// @Param retrievefornotifications body dto.RetrieveForNotificationsReq true "Retrieve"
// @Success 200 {object} string "recipients"
// @Failure 400 {string} string "invalid request"
// @Router /api/retrievefornotifications [post]
func RetrieveForNotifications(c *fiber.Ctx) error {
	var registeredStudentEmails []string
	var validStudentsFromNotification []string

	s := new(dto.RetrieveForNotificationsReq)
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
