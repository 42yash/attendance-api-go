package main

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Username string
	Name     string
	Claims   []ClaimReview `gorm:"foreignKey:TeacherId"`
}

func createTeacherHandler(w http.ResponseWriter, r *http.Request) {
	teacher := &Teacher{}
	err := json.NewDecoder(r.Body).Decode(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	result := db.Create(teacher)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
}

func getTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.URL.Query().Get("id")

	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	var teacher Teacher
	result := db.Preload("Claims").First(&teacher, teacherId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(teacher)
}

func getClaimsByTeacherIdHandler(w http.ResponseWriter, r *http.Request) {
	teacher := &Teacher{}
	err := json.NewDecoder(r.Body).Decode(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	result := db.Preload("Claims").First(&teacher, username)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(teacher.Claims)
}
