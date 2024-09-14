package file_service

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type FileMetadata struct {
	ID         int
	UserID     int
	Filename   string
	FilePath   string
	Size       int64
	UploadDate time.Time
}

func UploadFileHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user").(map[string]interface{})["user_id"].(int)

		file, handler, err := r.FormFile("file")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer file.Close()

		// saving file locally
		filepath := filepath.Join("uploads", handler.Filename)

		dst, err := os.Create(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		size, err := io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// sving data to db

		_, err = db.Exec("INSERT INTO files (user_id, filename, filepath, size, upload_date) VALUES ($1, $2, $3, $4, $5)",
			userID, handler.Filename, filepath, size, time.Now())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("File uploaded successfully"))
	}
}

func GetFilesHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user").(map[string]interface{})["user_id"].(int)

		rows, err := db.Query("SELECT id, filename, filepath, size, upload_date FROM files WHERE user_id = $1", userID)

		if err != nil {
			http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var files []FileMetadata

		for rows.Next() {
			var file FileMetadata
			rows.Scan(&file.ID, &file.Filename, &file.FilePath, &file.Size, &file.UploadDate)
			files = append(files, file)
		}

		json.NewEncoder(w).Encode(files)

	}
}

func ShareFileHandler(db *sql.DB, redis_client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// getting params
		fileID := mux.Vars(r)["file_id"]

		var file FileMetadata

		err := db.QueryRow("SELECT id, filename, filepath FROM files WHERE id = $1", fileID).Scan(&file.ID, &file.Filename, &file.FilePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// serving th file
		http.ServeFile(w, r, file.FilePath)
	}
}
