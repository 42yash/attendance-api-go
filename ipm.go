package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type IPM struct {
	gorm.Model
	Username string
	Name     string
}

func getAllClaims(w http.ResponseWriter, r *http.Request) {
	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	var claims []MedicalClaim
	// Find all Claims where none of the ClaimReview has a status of "pending"
	result := db.Preload("Student").Preload("ClaimReviews").Preload("ClaimReviews.Teacher").Preload("Files").Find(&claims, "id NOT IN (SELECT claim_id FROM claim_reviews WHERE status = 'pending')")
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(claims)
}

func updateClaimStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	claimId, err := strconv.Atoi(vars["claimid"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Parse JSON
	var medicalclaim MedicalClaim
	err = json.NewDecoder(r.Body).Decode(&medicalclaim)
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

	// Find the claim
	var claim MedicalClaim
	result := db.First(&claim, claimId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	// Update the claim
	result = db.Model(&claim).Updates(medicalclaim)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(claim)
}
