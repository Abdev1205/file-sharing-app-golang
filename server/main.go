package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Abdev1205/21BCE11045_Backend/internal/adapters/cache"
	"github.com/Abdev1205/21BCE11045_Backend/internal/adapters/database"
	"github.com/Abdev1205/21BCE11045_Backend/internal/application/auth_service"
	"github.com/Abdev1205/21BCE11045_Backend/internal/application/file_service"
	"github.com/Abdev1205/21BCE11045_Backend/pkg/config"
	"github.com/Abdev1205/21BCE11045_Backend/pkg/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// loading the env
	config.LoadConfig()

	// postgres database initialization
	db := database.ConnectPostgres()

	// redis cache initialization
	redisClient := cache.ConnectRedis()

	// Creating and ensuring the 'uploads' folder exists
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create uploads directory: %v", err)
		}
		log.Printf("Uploads directory created: %s", uploadDir)
	}

	// intialising the router
	router := mux.NewRouter()

	// Starting route
	// here i just want to say that welcome to abhay file sharing app backend

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Abhay File Sharing App Backend"))
	})

	// public routres taht we can show the content without valid jwt credentials
	router.HandleFunc("/register", auth_service.RegisterHandler(db)).Methods("POST")
	router.HandleFunc("/login", auth_service.LoginHandler(db)).Methods("POST")

	// protected routres taht we can show the content with valid jwt credentials only
	// so we will create a new sub route to provide the middleware support
	protected := router.NewRoute().Subrouter()

	// adding middleware so taht only autheticated user can access
	protected.Use(middleware.JWTAuthMiddleware)

	protected.HandleFunc("/upload", file_service.UploadFileHandler(db, redisClient)).Methods("POST")
	protected.HandleFunc("/files", file_service.GetFilesHandler(db, redisClient)).Methods("GET")
	protected.HandleFunc("/share/{id:[0-9]+}", file_service.ShareFileHandler(db, redisClient)).Methods("GET")
	protected.HandleFunc("/files/search", file_service.SearchFilesHandler(db, redisClient)).Methods("GET")

	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	// creating a server on port 8080
	log.Printf("Starting Server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
