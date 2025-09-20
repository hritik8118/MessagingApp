package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"MessagingApp/Backend/auth"
	"MessagingApp/Backend/db"
	"MessagingApp/Backend/models"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if db.UserExists(u.Username) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	u.CreatedAt = time.Now()
	db.SaveUser(u)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, ok := db.GetUser(creds.Username)
	if !ok || user.Password != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateJWT(user.Username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
