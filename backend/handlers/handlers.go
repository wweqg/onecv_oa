package handlers

import (
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/wweqg/onecv_oa/backend/database"
	"github.com/wweqg/onecv_oa/backend/models"
)

func ListTeachers(c *fiber.Ctx) error {
	teachers := []models.Teacher{}
	database.DB.Db.Find(&teachers)

	return c.Status(http.StatusOK).JSON(teachers)
}

func ListStudents(c *fiber.Ctx) error {
	students := []models.Student{}
	database.DB.Db.Find(&students)

	return c.Status(http.StatusOK).JSON(students)
}

func CreateTeacher(c *fiber.Ctx) error {
	var newTeacher models.Teacher

	if err := c.BodyParser(&newTeacher); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	if err := database.DB.Db.Create(&newTeacher).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "This teacher is already created",
		})
	}

	return c.Status(http.StatusOK).JSON(newTeacher)
}

func CreateStudent(c *fiber.Ctx) error {
	var newStudent struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&newStudent); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	if err := database.DB.Db.Create(&newStudent).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "This student is already created",
		})
	}

	return c.Status(http.StatusCreated).JSON(newStudent)
}

func RegisterStudentsToTeacher(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	// Find the teacher by email
	var teacher models.Teacher
	if err := db.Where("email = ?", request.Teacher).First(&teacher).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Teacher not found",
		})
	}

	// Register students to the teacher
	for _, studentEmail := range request.Students {
		var student models.Student
		if err := db.Where("email = ?", studentEmail).First(&student).Error; err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Student not found",
		})
		}

		if err := db.Create(&models.TeacherStudent{TeacherID: teacher.ID, StudentID: student.ID}).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "This student is already registered with this teacher",
			})
		}
	}

	return c.SendStatus(http.StatusNoContent)
}