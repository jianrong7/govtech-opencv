package controller

import (
	"govtech-opencv/db"
	"govtech-opencv/model"

	"github.com/gofiber/fiber/v2"
)

type RegisterStruct struct {
	Students []string `json:"students"`
	Teacher  string   `json:"teacher"`
}

func Register(c *fiber.Ctx) error {
	s := new(RegisterStruct)
	if err := c.BodyParser(s); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Register", "data": nil})
}

func GetCommonStudents(c *fiber.Ctx) error {
	teacherEmail := c.Query("teacher")
	var students []model.Student
	db.DB.Where("email = ?", teacherEmail).Find(&students)
	return c.JSON(fiber.Map{"students": students})
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
