package main

import (
	"encoding/json"
	"net/http"
)

func getStudentInfoHandler(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	if username == "" {
		http.Error(w, "Access denied", http.StatusUnauthorized)
		return
	}

	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	var student Student
	if result := db.Where("username = ?", username).First(&student); result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(student)
}
