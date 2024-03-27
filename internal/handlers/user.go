// internal/handlers/users.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shahin-bayat/scraper-api/internal/models"
)

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	// Decode request body into CreateUserRequest struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate age (assuming calculateAge function is defined elsewhere)
	age := calculateAge(req.DateOfBirth)

	// Create User object for database interaction
	user := models.User{
		Name: req.Name,
		Age:  age,
		// set other fields relevant to the database schema
	}

	// Use the User object for further processing...
	// For example, call UserRepository to create user
	err = h.store.UserRepository().CreateUser(r.Context(), &user)
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Convert User to CreateUserResponse
	resp := models.CreateUserResponse(user)

	// Respond with created user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func calculateAge(dateOfBirth string) int {
	// Implement age calculation
	return 0
}
