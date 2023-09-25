package handlers

import (
	"errors"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/wweqg/onecv_oa/backend/database"
	"github.com/wweqg/onecv_oa/backend/models"
	"gorm.io/gorm"
)

func ListStudents(c *fiber.Ctx) error {
	students := []models.Student{}
	database.DB.Db.Find(&students)

	return c.Status(http.StatusOK).JSON(students)
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

	// Check if the email field is empty
	if newStudent.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Email is required",
		})
	}

	// Check if the student with the same email already exists
	var existingStudent models.Student
	if err := database.DB.Db.Where("email = ?", newStudent.Email).First(&existingStudent).Error; err == nil {
		// Student with the same email already exists
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Student with this email already exists",
		})
	}

	// Create the new student in the database
	if err := database.DB.Db.Create(&models.Student{Email: newStudent.Email}).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create student",
		})
	}

	return c.Status(http.StatusCreated).JSON(newStudent)
}

func DeleteStudent(c *fiber.Ctx) error {
    db := database.DB.Db

    studentEmail := c.Params("email")

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

    if err := db.Delete(&student).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to delete student",
        })
    }

    return c.SendStatus(http.StatusNoContent)
}