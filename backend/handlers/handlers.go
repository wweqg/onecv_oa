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

func GetCommonStudents(c *fiber.Ctx) error {
	db := database.DB.Db

	// Get the teacher emails from the query parameters
	teacherEmails := c.Query("teacher")

	// Ensure at least one teacher email is provided
	if len(teacherEmails) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'teacher' query parameter",
		})
	}

	var commonStudents []models.Student
	db.Table("students").
		Joins("JOIN teacher_students ON students.id = teacher_students.student_id").
		Joins("JOIN teachers ON teacher_students.teacher_id = teachers.id").
		Where("teachers.email IN (?)", teacherEmails).
		Group("students.id").
		Find(&commonStudents)

	// Extract the common student emails
	commonStudentEmails := []string{}
	for _, student := range commonStudents {
		commonStudentEmails = append(commonStudentEmails, student.Email)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"students": commonStudentEmails,
	})
}

func SuspendStudent(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		Student string `json:"student"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	var student models.Student
	if err := db.Where("email = ?", request.Student).First(&student).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Student not found",
		})
	}

	// Suspend the student
	db.Model(&student).Update("suspend", true)

	return c.SendStatus(http.StatusNoContent)
}

// // RetrieveStudentsForNotifications retrieves a list of students who can receive a given notification.
// func RetrieveStudentsForNotifications(c *fiber.Ctx) error {
// 	db := c.Locals("db").(*gorm.DB)

// 	// Define request structure
// 	var request struct {
// 		Teacher      string `json:"teacher"`
// 		Notification string `json:"notification"`
// 	}

// 	// Parse the request body into the request struct
// 	if err := c.BodyParser(&request); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
// 	}

// 	// Find the teacher by email
// 	var teacher models.Teacher
// 	if err := db.Where("email = ?", request.Teacher).First(&teacher).Error; err != nil {
// 		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Teacher not found"})
// 	}

// 	// Parse the notification text to extract mentioned student emails
// 	mentionedStudents := extractMentionedStudents(request.Notification)

// 	// Find the students who can receive the notification
// 	var recipients []models.Student
// 	db.Table("students").
// 		Joins("JOIN teacher_students ON students.id = teacher_students.student_id").
// 		Joins("JOIN teachers ON teacher_students.teacher_id = teachers.id").
// 		Where("teachers.email = ?", teacher.Email).
// 		Where("students.suspend = ?", false).
// 		Where("students.email IN (?)", mentionedStudents).
// 		Find(&recipients)

// 	// Extract the recipient student emails
// 	recipientEmails := []string{}
// 	for _, recipient := range recipients {
// 		recipientEmails = append(recipientEmails, recipient.Email)
// 	}

// 	return c.JSON(fiber.Map{"recipients": recipientEmails})
// }

// // Helper function to extract mentioned student emails from notification text
// func extractMentionedStudents(notificationText string) []string {
// 	// Implement your logic to extract mentioned student emails here
// 	// For example, you can use regular expressions to find email addresses prefixed with '@'
// 	// and extract them into a slice.
// 	// Return the list of mentioned student emails.
// 	return []string{} // Replace with the actual implementation
// }
