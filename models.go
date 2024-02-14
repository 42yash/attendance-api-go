package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        // Includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"uniqueIndex"` // Ensures usernames are unique
	Password   string `gorm:"size:60"`
	UserType   string
}

type Student struct {
	gorm.Model                   // Includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	UserID               uint    // Foreign key for the User
	Name                 string  // Student's full name
	RegisterNumber       string  // Unique registration number for the student
	Class                string  // Class or course the student is enrolled in
	Email                string  // Student's email address
	AttendancePercentage float64 // Student's attendance percentage
}
