package models

import "gorm.io/gorm"

type Teacher struct {
    gorm.Model
    Email string `gorm:"uniqueIndex"`
}

type Student struct {
    gorm.Model
    Email   string `gorm:"uniqueIndex"`
    Suspend bool 
}

type TeacherStudent struct {
    gorm.Model
    TeacherID uint `gorm:"uniqueIndex:idx_teacher_student"`
    StudentID uint `gorm:"uniqueIndex:idx_teacher_student"`
}

type Notification struct {
    gorm.Model
    TeacherID   uint
    Text        string
}