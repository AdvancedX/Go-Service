package handlers

import (
	"Go-Service/database"
	"Go-Service/middleware"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var newPost database.Post
	newPost.AuthorID, _ = middleware.GetUserID(r.Context())
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "Failed to parse request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.CreatePost(&newPost); err != nil {
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post created successfully",
	})
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	idStr = strings.TrimSuffix(idStr, "/")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	existing, err := database.GetPostByID(uint64(postID)) // 需要你实现：可查 draft/published
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if existing.AuthorID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := database.DeletePostByID(uint64(postID)); err != nil {
		http.Error(w, "Failed to delete post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post deleted successfully",
	})
}
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	idStr = strings.TrimSuffix(idStr, "/")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	existing, err := database.GetPostByID(uint64(postID)) // 需要你实现：可查 draft/published
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if existing.AuthorID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Title     string `json:"title"`
		ContentMD string `json:"content_md"`
		Status    string `json:"status"` // "draft" or "published"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.ContentMD) == "" {
		http.Error(w, "title and content_md are required", http.StatusBadRequest)
		return
	}
	if req.Status != "" && req.Status != "draft" && req.Status != "published" {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	existing.Title = req.Title
	existing.ContentMD = req.ContentMD
	if req.Status != "" {
		existing.Status = req.Status
	}

	if err := database.UpdatePost(existing); err != nil {
		http.Error(w, "Failed to update post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(existing)
}
func ListPublishedPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	posts, err := database.ListPublishedPosts(1, 10)
	if err != nil {
		http.Error(w, "Failed to list posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
func GetPublishedPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input ID", http.StatusBadRequest)
		return
	}
	post, err := database.GetPublishedPostByID(uint64(id))
	if err != nil {
		http.Error(w, "Failed to get post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
