package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func initDB() {
	db, sqlDB, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Perform the migration
	if err := db.AutoMigrate(&User{}, &Student{}); err != nil {
		log.Fatalf("Error auto migrating tables: %v", err)
	}

	// Print Success
	fmt.Println("Database initialization successful.")
}

func initServer() {
	router := mux.NewRouter()

	// Apply the middleware to the router
	router.Use(verifyTokenMiddleware)

	// Register unprotected routes
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/register", registerHandler).Methods("POST")

	// Register protected routes
	studentRouter := router.PathPrefix("/student").Subrouter()
	studentRouter.Use(authorizeRole("student"))
	studentRouter.HandleFunc("/info", getStudentInfo).Methods("GET")

	// ... register other routes for Teacher and IPM ...

	// Print the message
	fmt.Println("Server starting on http://localhost:8000...")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	initDB()
	initServer()
}
