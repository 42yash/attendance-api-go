package main

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	// ... other imports
)

type Student struct {
	gorm.Model                        // Includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	Username             string       // Foreign key for the User
	Name                 string       // Student's full name
	Class                string       // Class or course the student is enrolled in
	RegisterNumber       string       // Unique registration number for the student
	Email                string       // Student's email address
	Phone                string       // Student's phone number
	AttendancePercentage float64      // Student's attendance percentage
	Attendance           []Attendance `gorm:"foreignKey:StudentId"`
}

func createStudentInfo(w http.ResponseWriter, r *http.Request) {
	// Extract the username from JWT claims
	username, err := getUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT", http.StatusInternalServerError)
		return
	}

	// Decode the request body into a Student struct
	var student Student
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Find the student and preload the attendance records
	result := db.Preload("Attendance").Where("username = ?", username).First(&student)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	// Respond with the student info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func getStudentInfo(w http.ResponseWriter, r *http.Request) {
	// Extract the username from JWT claims
	username, err := getUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Fetch student info from the database and preload the attendance records
	var student Student
	result := db.Preload("Attendance").Where("username = ?", username).First(&student)
	if result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Return student info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}
