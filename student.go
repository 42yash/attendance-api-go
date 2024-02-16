package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// ... other imports
)

func test(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT", http.StatusInternalServerError)
		return
	}
	fmt.Println(username)
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

	// Create a new student record in the database
	student.Username = username
	result := db.Create(&student)
	if result.Error != nil {
		http.Error(w, "Failed to create student record", http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.WriteHeader(http.StatusCreated)

}

func getStudentInfo(w http.ResponseWriter, r *http.Request) {
	// Extract the username from JWT claims
	username, err := getUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT", http.StatusInternalServerError)
		return
	}
	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Fetch student info from the database
	var student Student
	result := db.Where("username = ?", username).First(&student)
	if result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	fmt.Println(student)

	// Return student info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}
