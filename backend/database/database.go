package database

import (
	"fmt"
	"os"

	"github.com/wweqg/onecv_oa/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Adapted from https://github.com/divrhino/divrhino-trivia/blob/main/database/database.go

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		os.Exit(2)
	}

	db.AutoMigrate(&models.Teacher{})
	db.AutoMigrate(&models.Student{})
	db.AutoMigrate(&models.TeacherStudent{})
	db.AutoMigrate(&models.Notification{})

	DB = Dbinstance{
		Db: db,
	}
}