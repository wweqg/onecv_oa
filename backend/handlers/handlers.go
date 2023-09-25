package handlers

import (
	"errors"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/wweqg/onecv_oa/backend/database"
	"github.com/wweqg/onecv_oa/backend/models"
	"gorm.io/gorm"
)

func ListTeachersStudents(c *fiber.Ctx) error {
	list := []models.TeacherStudent{}
	database.DB.Db.Find(&list)

	return c.Status(http.StatusOK).JSON(list)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Teacher not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch teacher",
		})
	}

	// Register students to the teacher
	for _, studentEmail := range request.Students {
		var student models.Student
		if err := db.Where("email = ?", studentEmail).First(&student).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"message": "Student not found",
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch student",
			})
		}

		// Check if the student is already registered with the teacher
		var existingRegistration models.TeacherStudent
		if err := db.Where("teacher_id = ? AND student_id = ?", teacher.ID, student.ID).First(&existingRegistration).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				// Error other than not found occurred, return an error
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to check existing registration",
				})
			}

			// Student is not registered with the teacher, so register them
			if err := db.Create(&models.TeacherStudent{TeacherID: teacher.ID, StudentID: student.ID}).Error; err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to register student with teacher",
				})
			}
		}
	}

	return c.SendStatus(http.StatusNoContent)
}
