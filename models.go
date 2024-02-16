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
	Username             string  // Foreign key for the User
	Name                 string  // Student's full name
	Class                string  // Class or course the student is enrolled in
	RegisterNumber       string  // Unique registration number for the student
	Email                string  // Student's email address
	Phone                string  // Student's phone number
	AttendancePercentage float64 // Student's attendance percentage
	Attendance           []Attendance
}

type Teacher struct {
	Username string
}

type Attendance struct {
	gorm.Model
	Student   Student
	Course    string
	Time      string
	Date      string
	Teacher   Teacher
	IsPresent bool
	Applied   bool
	Claimed   bool
}
