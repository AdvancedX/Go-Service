package main

import (
	"Go-Service/database"
	"Go-Service/handlers"
	"Go-Service/middleware"
	"fmt"
	"net/http"
)

func main() {
	database.ConnectDatabase()

	// HTTP 路由
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/auth/register", handlers.RegisterHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)
	http.Handle("/users", middleware.JWTAuth(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			handlers.CreateUserHandler(writer, request)
		case http.MethodGet:
			handlers.ListUsersHandler(writer, request)
		default:
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	http.Handle("/users/", middleware.JWTAuth(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handlers.FindUserByIDHandler(writer, request)
		case http.MethodDelete:
			handlers.DeleteUserByIDHandler(writer, request)
		case http.MethodPut:
			handlers.UpdateUserHandler(writer, request)
		default:
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	// 启动 HTTP 服务，监听端口 8080
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
