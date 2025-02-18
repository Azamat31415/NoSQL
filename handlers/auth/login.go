package auth

import (
	"GoProject/configs"
	"GoProject/migrations"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
}

func GenerateJWT(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.JwtSecret)
}

func LoginHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginData LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var user migrations.User
		err := collection.FindOne(context.TODO(), bson.M{"email": loginData.Email}).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if !user.CheckPassword(loginData.Password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := GenerateJWT(user.ID.Hex(), user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		response := LoginResponse{
			Message: "Login successful",
			Token:   token,
			UserID:  user.ID.Hex(),
			Role:    user.Role,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
