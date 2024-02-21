package main

import "gorm.io/gorm"

type MedicalClaim struct {
	gorm.Model
	StudentId    uint
	Reason       string
	Description  string
	Status       string
	ClaimReviews []ClaimReview `gorm:"foreignKey:ClaimId"`
}

type ClaimReview struct {
	gorm.Model
	ClaimId      uint         // Foreign key to the MedicalClaim
	MedicalClaim MedicalClaim `gorm:"foreignKey:ClaimId"`
	AttendanceId uint         // Foreign key to the Attendance
	Attendance   Attendance   `gorm:"foreignKey:AttendanceId"`
	TeacherId    string       // Foreign key to the Teacher
	Approved     bool         // Whether the claim was approved or rejected
	Message      string       // Optional message left by the teacher
}
