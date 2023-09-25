package main

import (
	"github.com/wweqg/onecv_oa/backend/database"
	"github.com/wweqg/onecv_oa/backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()

    app := fiber.New()

	setupRoutes(app)

    app.Listen(":3000")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/teachers", handlers.ListTeachers)
	app.Get("/api/students", handlers.ListStudents)

	app.Post("/api/teachers", handlers.CreateTeacher)
	app.Post("/api/students", handlers.CreateStudent)

	app.Post("/api/register", handlers.RegisterStudentsToTeacher)
}