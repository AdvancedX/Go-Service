package middleware

import (
	"Go-Service/auth"
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 JWT token
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ParseJWT(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int64, bool) {
	v := ctx.Value(UserIDKey)
	id, ok := v.(int64)
	return id, ok
}
