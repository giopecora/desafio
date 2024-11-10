package handlers

import (
	"encoding/json"
	"main/auth"
	"main/services"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(client *mongo.Client) *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(client),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		responseJSON(w, http.StatusBadRequest, AppError{Message: "Invalid request body", Code: http.StatusBadRequest})
		return
	}

	user, err := h.userService.Authenticate(loginReq.Username, loginReq.Password)
	if err != nil {
		responseJSON(w, http.StatusUnauthorized, AppError{Message: "Invalid credentials", Code: http.StatusUnauthorized})
		return
	}

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, AppError{Message: ErrInternalServer, Code: http.StatusInternalServerError})
		return
	}
	token, err := auth.GenerateToken(user.ID.Hex(), user.IsAdmin)
	if err != nil {
		responseJSON(w, http.StatusUnauthorized, AppError{Message: "Invalid credentials", Code: http.StatusUnauthorized})
		return
	}
	responseJSON(w, http.StatusOK, map[string]string{"token": token})
}
