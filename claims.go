package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name           string
	Path           string
	MedicalClaimID uint
}
type MedicalClaim struct {
	gorm.Model
	StudentId    uint
	Student      Student `gorm:"foreignKey:StudentId"`
	Reason       string
	Description  string
	Status       string        `gorm:"default:Pending"`
	ClaimReviews []ClaimReview `gorm:"foreignKey:ClaimId"`
	Files        []File        `gorm:"foreignKey:MedicalClaimID"`
}

type ClaimReview struct {
	gorm.Model
	ClaimId      uint         // Foreign key to the MedicalClaim
	MedicalClaim MedicalClaim `gorm:"foreignKey:ClaimId"`
	AttendanceId uint         // Foreign key to the Attendance
	Attendance   Attendance   `gorm:"foreignKey:AttendanceId"`
	TeacherId    string       // Foreign key to the Teacher
	Teacher      Teacher      `gorm:"foreignKey:TeacherId"`
	Status       string       `gorm:"default:Pending"`
	Message      string       // Optional message left by the teacher
}

type RequestBody struct {
	Reason      string   `json:"reason"`
	Description string   `json:"description"`
	Data        []string `json:"data"`
	Date        []string
	Period      []string
	Files       []string `json:"files"`
	FileNames   []string `json:"filenames"`
}

func getAttendanceRecords(date string, period string, studentId uint) ([]Attendance, error) {
	db, sqlDB, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	var records []Attendance
	result := db.Where("date = ? AND period = ? AND student_id = ?", date, period, studentId).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func createMedicalClaim(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set Reason and Description
	var medicalClaim MedicalClaim
	medicalClaim.Reason = requestBody.Reason
	medicalClaim.Description = requestBody.Description

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

	// Get studentid from username
	var student Student
	result := db.Preload("Attendance").Where("username = ?", username).First(&student)
	if result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	medicalClaim.StudentId = student.ID

	// Save medicalClaim to the database
	result = db.Create(&medicalClaim)
	if result.Error != nil {
		http.Error(w, "Failed to save medical claim", http.StatusInternalServerError)
		return
	}

	if len(requestBody.Files) > 0 {
		for i, file := range requestBody.Files {
			fileRecord := File{
				Path:           file,
				Name:           requestBody.FileNames[i],
				MedicalClaimID: medicalClaim.ID, // Link the file to the medical claim
			}

			result := db.Create(&fileRecord)
			if result.Error != nil {
				http.Error(w, "Failed to save file record", http.StatusInternalServerError)
				return
			}
		}
	}

	// Fetch all attendance using requestBody from db

	for _, dp := range requestBody.Data {
		period := dp[len(dp)-2:]
		date := dp[:len(dp)-3]

		attendanceRecords, err := getAttendanceRecords(date, period, student.ID)
		if err != nil {
			http.Error(w, "Failed to fetch attendance records", http.StatusInternalServerError)
			return
		}

		for _, attendanceRecord := range attendanceRecords {
			claimReview := ClaimReview{
				ClaimId:      medicalClaim.ID,
				AttendanceId: attendanceRecord.ID,
				TeacherId:    attendanceRecord.TeacherId,
				Status:       "pending",
			}

			result := db.Create(&claimReview)
			if result.Error != nil {
				http.Error(w, "Failed to create claim review", http.StatusInternalServerError)
				return
			}
		}
	}

	// Respond with newly created medical claim
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medicalClaim)
}

func getMedicalClaimByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	claimId, err := strconv.Atoi(vars["claimid"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, sqlDB, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	var medicalClaim MedicalClaim
	result := db.Preload("Student").Preload("ClaimReviews").Preload("Files").Preload("ClaimReviews.Teacher").Preload("ClaimReviews.Attendance").Where("id = ?", claimId).First(&medicalClaim)
	if result.Error != nil {
		http.Error(w, "Medical claim not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(medicalClaim)
}

func getClaimsByStudentHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch all medical claims for the student
	var student Student
	result := db.Preload("MedicalClaims").Where("username = ?", username).First(&student)
	if result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(student.MedicalClaims)
}

func getClaimsByTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch all medical claims for the student
	var teacher Teacher
	result := db.Preload("Claim").Preload("Claim.MedicalClaim").Preload("Claim.MedicalClaim.Files").Preload("Claim.MedicalClaim.Student").Where("username = ?", username).First(&teacher)
	if result.Error != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(teacher.Claim)
}
