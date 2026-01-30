package main

import (
	"Go-Service/handlers"
	"fmt"
	"net/http"
)

func main() {
	// HTTP 路由
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			handlers.CreateUserHandler(writer, request)
		case http.MethodGet:
			if id := request.URL.Query().Get("id"); id == "" {
				handlers.ListUsersHandler(writer, request)
			} else {
				handlers.FindUserByIDHandler(writer, request)
			}
		case http.MethodDelete:
			handlers.DeleteUserByIDHandler(writer, request)
		case http.MethodPut:
			handlers.UpdateUserHandler(writer, request)
		}
	})
	// 启动 HTTP 服务，监听端口 8080
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
