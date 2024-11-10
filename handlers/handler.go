package handlers

import (
	"encoding/json"
	"main/auth"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type AppError struct {
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Message
}

const (
	ErrPermissionDenied = "Permission denied"
	ErrInvalidData      = "Invalid data"
	ErrInvalidID        = "Invalid ID"
	ErrInternalServer   = "Internal server error"
)

func responseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			responseJSON(w, http.StatusUnauthorized, AppError{Message: "Missing authorization header", Code: http.StatusUnauthorized})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			responseJSON(w, http.StatusUnauthorized, AppError{Message: "Invalid authorization header", Code: http.StatusUnauthorized})
			return
		}

		claims, err := auth.ValidateToken(bearerToken[1])
		if err != nil {
			responseJSON(w, http.StatusUnauthorized, AppError{Message: "Invalid token", Code: http.StatusUnauthorized})
			return
		}
		r = r.WithContext(setUserContext(r.Context(), claims.UserID, claims.IsAdmin))

		next.ServeHTTP(w, r)
	}
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if !isAdmin(r) {
			responseJSON(w, http.StatusForbidden, AppError{Message: ErrPermissionDenied, Code: http.StatusForbidden})
			return
		}
		next.ServeHTTP(w, r)
	})
}
func UserMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if isAdmin(r) {
			responseJSON(w, http.StatusForbidden, AppError{Message: ErrPermissionDenied, Code: http.StatusForbidden})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAdmin(r *http.Request) bool {
	_, isAdmin := getUserFromContext(r.Context())
	return isAdmin
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Second), 1000)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
