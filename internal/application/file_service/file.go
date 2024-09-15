package file_service

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type FileMetadata struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Filename   string    `json:"filename"`
	FilePath   string    `json:"file_path"`
	Size       int64     `json:"size"`
	UploadDate time.Time `json:"upload_date"`
}

/// I know this file gone too messy but I don't have much time rectiy it

// Helper function to extract the user ID from JWT claims
func getUserIDFromContext(r *http.Request) (int, error) {
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		return 0, http.ErrNoCookie
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, http.ErrNoCookie
	}

	return int(userIDFloat), nil
}

func UploadFileHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UploadFileHandler called")

		// Extracting user ID from context
		userID, err := getUserIDFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieving the file from the request /// most painful part
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Ensuring the uploads directory exists // checking
		uploadDir := "uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err := os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
				return
			}
		}

		// Saving the file to the uploads directory
		filePath := filepath.Join(uploadDir, handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copying the file content to the destination
		size, err := io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Saving file metadata to the database
		_, err = db.Exec(
			"INSERT INTO files (user_id, filename, filepath, size, upload_date) VALUES ($1, $2, $3, $4, $5)",
			userID, handler.Filename, filePath, size, time.Now(),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileURL := "http://" + r.Host + "/uploads/" + handler.Filename

		// Prepareing JSON response
		response := map[string]string{
			"message":  "File uploaded successfully",
			"file_url": fileURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send JSON response
		json.NewEncoder(w).Encode(response)
	}
}

func GetFilesHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context
		userID, err := getUserIDFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query("SELECT id, filename, filepath, size, upload_date FROM files WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parsing the rows into a slice of FileMetadata structs
		var files []FileMetadata
		for rows.Next() {
			var file FileMetadata
			err := rows.Scan(&file.ID, &file.Filename, &file.FilePath, &file.Size, &file.UploadDate)
			if err != nil {
				http.Error(w, "Failed to scan file data", http.StatusInternalServerError)
				return
			}
			files = append(files, file)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(files)
	}
}

func ShareFileHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileID := vars["id"]
		log.Printf("File ID: %s", fileID)

		if fileID == "" {
			log.Println("Error: File ID is missing from URL")
			http.Error(w, "File ID is missing", http.StatusBadRequest)
			return
		}

		var file FileMetadata
		err := db.QueryRow("SELECT id, filename, filepath FROM files WHERE id = $1", fileID).Scan(&file.ID, &file.Filename, &file.FilePath)
		if err != nil {
			log.Printf("Database error or file not found: %v", err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		println("file metadata", file.FilePath, file.Filename, file.ID)
		fileURL := "http://" + r.Host + "/uploads/" + file.Filename

		response := map[string]string{
			"message":  "File shared successfully",
			"file_url": fileURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)

	}
}

func SearchFilesHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Getting query parameters for search and tehy can be empty as well so it will get all the files
		queryName := r.URL.Query().Get("name")
		queryType := r.URL.Query().Get("type")
		queryDate := r.URL.Query().Get("date") // Expected in YYYY-MM-DD format

		// dynamic sql query
		query := "SELECT id, filename, filepath, size, upload_date FROM files WHERE user_id = $1"
		params := []interface{}{userID}

		// condition // can be further optimised but
		if queryName != "" {
			query += " AND LOWER(filename) LIKE $2"
			params = append(params, "%"+strings.ToLower(queryName)+"%")
		}
		if queryType != "" {
			query += " AND LOWER(filepath) LIKE $3"
			params = append(params, "%"+strings.ToLower(queryType)+"%")
		}
		if queryDate != "" {
			query += " AND DATE(upload_date) = $4"
			parsedDate, err := time.Parse("2006-01-02", queryDate)
			if err != nil {
				http.Error(w, "Invalid date format. Expected YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			params = append(params, parsedDate)
		}

		// passing the quqry and params for serched content
		rows, err := db.Query(query, params...)
		if err != nil {
			log.Printf("Failed to execute search query: %v", err)
			http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// json stuff
		var files []FileMetadata
		for rows.Next() {
			var file FileMetadata
			err := rows.Scan(&file.ID, &file.Filename, &file.FilePath, &file.Size, &file.UploadDate)
			if err != nil {
				log.Printf("Failed to scan file data: %v", err)
				http.Error(w, "Failed to scan file data", http.StatusInternalServerError)
				return
			}
			files = append(files, file)
		}

		// Return the search results as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(files)
	}
}
