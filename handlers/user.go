package handlers

import (
	"Go-Service/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// CreateUserHandler Create User
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var newUser database.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Failed to parse request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.CreateUser(&newUser); err != nil {
		// 如果数据库写入失败，返回错误
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 响应成功状态
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

// DeleteUserByIDHandler Delete User By ID
func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input ID", http.StatusBadRequest)
		return
	}
	if err := database.DeleteUserByID(id); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})

}

// ListUsersHandler List
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input ID", http.StatusBadRequest)
		return
	}
	var updatedUser database.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Failed to parse request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.UpdateUser(id, &updatedUser); err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

func FindUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	var id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input ID", http.StatusBadRequest)
		return
	}
	users, err := database.GetUserByID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
