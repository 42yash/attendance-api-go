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
	if err := db.AutoMigrate(&User{}, &Student{}); err != nil {
		log.Fatalf("Error auto migrating tables: %v", err)
	}

	// Print Success
	fmt.Println("Database initialization successful.")
}

func initServer() {
	router := mux.NewRouter()

	// Register unprotected routes
	router.HandleFunc("/login", loginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", registerHandler).Methods("POST")

	// Register protected routes
	studentRouter := router.PathPrefix("/student").Subrouter()
	studentRouter.Use(authorizeRole("student"))
	studentRouter.HandleFunc("/info", getStudentInfo).Methods("GET")
	studentRouter.HandleFunc("/create", createStudentInfo).Methods("POST")
	studentRouter.HandleFunc("/test", test).Methods("GET")

	// Apply other middleware to the router
	router.Use(jsonContentTypeMiddleware)

	// ... register other routes for Teacher and IPM ...

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Print the message
	fmt.Println("Server starting on http://localhost:8000...")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func main() {
	initDB()
	initServer()
}
