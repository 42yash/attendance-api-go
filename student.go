package main

import (
	"encoding/json"
	"net/http"
	// ... other imports
)

func getStudentInfo(w http.ResponseWriter, r *http.Request) {
	// Extract the username from JWT claims
	claims := r.Context().Value("claims").(*CustomClaims)

	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Fetch student info from the database
	var student Student
	result := db.Where("username = ?", claims.Username).First(&student)
	if result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Return student info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}
