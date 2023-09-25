package handlers

import (
	"errors"
	"net/http"
	"regexp"
    "strings"
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

	db.Model(&student).Update("suspend", true)

	return c.SendStatus(http.StatusNoContent)
}

func RetrieveStudentsForNotifications(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	var teacher models.Teacher
	if err := db.Where("email = ?", request.Teacher).First(&teacher).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Teacher not found",
		})
	}

	mentionedStudents := extractMentionedStudents(request.Notification)

	var recipient_1 []models.Student

	var recipient_2 []models.Student

	var recipient_3 []models.Student

	// Subquery 2: Students registered with this teacher
	db.Table("students").
		Joins("JOIN teacher_students ON students.id = teacher_students.student_id").
		Joins("JOIN teachers ON teacher_students.teacher_id = teachers.id").
		Where("teachers.email = ?", teacher.Email).
		Find(&recipient_1)

	// Subquery 2: Students mentioned in the notification
	db.Table("students").
		Where("students.email IN (?)", mentionedStudents).
		Find(&recipient_2)

	// Subquery 3: Students not suspended
	db.Table("students").
		Where("students.suspend = ?", false).
		Find(&recipient_3)

	recipients := unionAndIntersect(recipient_1, recipient_2, recipient_3)

	recipientEmails := []string{}
	for _, recipient := range recipients {
		recipientEmails = append(recipientEmails, recipient.Email)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"recipients": recipientEmails,
	})
}

func unionAndIntersect(s1, s2, s3 []models.Student) []models.Student {
    // Convert slices to sets
    set1 := make(map[string]bool)
    set2 := make(map[string]bool)
    set3 := make(map[string]bool)

    for _, student := range s1 {
        set1[student.Email] = true
    }

    for _, student := range s2 {
        set2[student.Email] = true
    }

    for _, student := range s3 {
        set3[student.Email] = true
    }

    // Perform union of set1 and set2
    unionSet := make(map[string]bool)
    for k := range set1 {
        unionSet[k] = true
    }
    for k := range set2 {
        unionSet[k] = true
    }

    // Perform intersection of unionSet and set3
    result := []models.Student{}
    for k := range unionSet {
        if set3[k] {
            result = append(result, models.Student{Email: k})
        }
    }

    return result
}


func extractMentionedStudents(notificationText string) []string {

	//regular expression to find email addresses prefixed with '@'
	mentionRegex := regexp.MustCompile(`@([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4})`)
	
	// Find all email mentions in the notification text
	matches := mentionRegex.FindAllString(notificationText, -1)
	
	// Extract email addresses from mentions (remove '@' symbol)
	mentionedStudents := []string{}
	for _, match := range matches {
		email := strings.TrimPrefix(match, "@")
		mentionedStudents = append(mentionedStudents, email)
	}
	
	return mentionedStudents
}
