package handlers

import (
	"errors"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/wweqg/onecv_oa/backend/database"
	"github.com/wweqg/onecv_oa/backend/models"
	"gorm.io/gorm"
)

func ListTeachers(c *fiber.Ctx) error {
	teachers := []models.Teacher{}
	database.DB.Db.Find(&teachers)

	return c.Status(http.StatusOK).JSON(teachers)
}

func CreateTeacher(c *fiber.Ctx) error {
	var newTeacher models.Teacher

	if err := c.BodyParser(&newTeacher); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	// Check if the email field is empty
	if newTeacher.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Email is required",
		})
	}

	// Check if the teacher with the same email already exists
	var existingTeacher models.Teacher
	if err := database.DB.Db.Where("email = ?", newTeacher.Email).First(&existingTeacher).Error; err == nil {
		// Teacher with the same email already exists
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Teacher with this email already exists",
		})
	}

	// Create the new teacher in the database
	if err := database.DB.Db.Create(&newTeacher).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create teacher",
		})
	}

	return c.Status(http.StatusCreated).JSON(newTeacher)
}

func DeleteTeacher(c *fiber.Ctx) error {
    db := database.DB.Db

    teacherEmail := c.Params("email")

    var teacher models.Teacher
	if err := db.Where("email = ?", teacherEmail).First(&teacher).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(http.StatusNotFound).JSON(fiber.Map{
                "message": "Teacher not found",
            })
        }
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to fetch teacher",
        })
    }

    if err := db.Delete(&teacher).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to delete teacher",
        })
    }

    return c.SendStatus(http.StatusNoContent)
}