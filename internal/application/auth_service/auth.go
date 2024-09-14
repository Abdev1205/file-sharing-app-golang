package auth_service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Abdev1205/21BCE11045_Backend/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to Hash Password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, hashedPassword)
		if err != nil {
			http.Error(w, "Failed to Register User", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		var storedPassword string

		// here first we are checking that user is already registered on not
		// if not registered so we are not going forward
		err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", user.Email).Scan(&user.ID, &storedPassword)

		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// now user is valid and verified so we have to assig them token
		// so that user can interact with us using that jwt token as identity

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(72 * time.Hour).Unix(),
		})

		tokenStr, err := token.SignedString([]byte(config.GetJWTSecret()))

		if err != nil {
			http.Error(w, "Failed to generate JWT Token", http.StatusInternalServerError)
			return
		}

		// now we have to this tokenstr to cookies in the frontend so that
		// when user interact so first we will check their cookies val
		// if cookies val is valid then only we will proceed

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenStr,
			Expires:  time.Now().Add(72 * time.Hour),
			HttpOnly: true,
			Secure:   false,
		})

		// now we will return reponse as user and message user Logged in Successfully
		// Send the response with user data and success message
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message": "User logged in successfully",
			"user": map[string]interface{}{
				"id":    user.ID,
				"email": user.Email,
			},
		}
		json.NewEncoder(w).Encode(response)

	}
}
