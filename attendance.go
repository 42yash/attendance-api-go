package main

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	StudentId uint
	Course    string
	Period    string
	Date      string
	TeacherId string
	IsPresent bool
	IsApplied bool
	IsClaimed bool
}

func createAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new 'Attendance' instance
	attendance := &Attendance{}
	err := json.NewDecoder(r.Body).Decode(attendance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Find the student
	var student Student
	result := db.Preload("Attendance").First(&student, attendance.StudentId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	// Insert the new attendance into the database
	result = db.Create(attendance)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Append the new attendance to the student's attendance slice
	student.Attendance = append(student.Attendance, *attendance)

	// Save the student
	result = db.Save(&student)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the newly created attendance
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendance)
}
