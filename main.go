package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	if err := db.AutoMigrate(&User{}, &Student{}, &Attendance{}, &MedicalClaim{}, &Teacher{}, &ClaimReview{}, &File{}, &IPM{}); err != nil {
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
	claimsRouter.HandleFunc("/{claimid}", getMedicalClaimByIdHandler).Methods("GET")
	claimsRouter.HandleFunc("/", getClaimsByStudentHandler).Methods("GET")

	// /teacher routes
	teacherRouter := router.PathPrefix("/teacher").Subrouter()
	teacherRouter.HandleFunc("/self", getTeacherByTokenHandler).Methods("GET")
	teacherRouter.HandleFunc("/create", createTeacherHandler).Methods("POST")
	teacherRouter.HandleFunc("/claims", getClaimsByTeacherHandler).Methods("GET")
	teacherRouter.HandleFunc("/claims/{claimid}", putClaimReviewHandler).Methods("PUT")

	// /ipm routes
	ipmRouter := router.PathPrefix("/ipm").Subrouter()
	ipmRouter.HandleFunc("/claims", getAllClaims).Methods("GET")
	ipmRouter.HandleFunc("/claims/{claimid}", updateClaimStatus).Methods("PUT")

	// Apply other middleware to the router
	router.Use(jsonContentTypeMiddleware)
	router.Use(loggingMiddleware)

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONTEND_URL")},
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
