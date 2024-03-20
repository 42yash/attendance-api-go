package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectDB() (*gorm.DB, *sql.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
