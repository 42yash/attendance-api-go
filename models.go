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

type Teacher struct {
	Username string
}
