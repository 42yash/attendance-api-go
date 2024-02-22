package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func initDB() {
	db, sqlDB, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Perform the migration
	if err := db.AutoMigrate(&User{}, &Student{}, &Attendance{}, &MedicalClaim{}, &ClaimReview{}); err != nil {
		log.Fatalf("Error auto migrating tables: %v", err)
	}

	// Print Success
	fmt.Println("Database initialization successful.")
}

func initServer() {
	router := mux.NewRouter()

	// Register unprotected routes
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/register", registerHandler).Methods("POST")

	// /student routes
	studentRouter := router.PathPrefix("/student").Subrouter()
	studentRouter.Use(authorizeRole("student"))
	studentRouter.HandleFunc("/create", createStudentInfo).Methods("POST")
	studentRouter.HandleFunc("/info", getStudentInfo).Methods("GET")

	// /attendance routes
	attendanceRouter := router.PathPrefix("/attendance").Subrouter()
	attendanceRouter.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		authorizeRole("admin")(http.HandlerFunc(createAttendanceHandler)).ServeHTTP(w, r)
	}).Methods("POST")

	// /claims routes
	claimsRouter := router.PathPrefix("/claims").Subrouter()
	claimsRouter.HandleFunc("/create", createMedicalClaim).Methods("POST")

	// Apply other middleware to the router
	router.Use(jsonContentTypeMiddleware)
	router.Use(loggingMiddleware)

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	handler := corsWrapper.Handler(router)

	// Print the message
	fmt.Println("Server starting on http://localhost:8000...")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func main() {
	initDB()
	initServer()
}
