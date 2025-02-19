package auth

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func RegisterHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var regData migrations.UserRegistration
		if err := json.NewDecoder(r.Body).Decode(&regData); err != nil {
			http.Error(w, "Invalid data format", http.StatusBadRequest)
			return
		}

		if regData.Email == "" || regData.Password == "" || regData.FirstName == "" {
			http.Error(w, "Fields email, password, and first_name are required", http.StatusBadRequest)
			return
		}

		// Проверка на существующего пользователя
		var existingUser migrations.User
		err := collection.FindOne(context.TODO(), bson.M{"email": regData.Email}).Decode(&existingUser)
		if err == nil {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}

		// Создание нового пользователя
		user := migrations.User{
			ID:        primitive.NewObjectID(),
			Email:     regData.Email,
			Password:  regData.Password,
			FirstName: regData.FirstName,
			LastName:  regData.LastName,
			Address:   regData.Address,
			Phone:     regData.Phone,
			Role:      "user",
		}

		if err := user.HashPassword(); err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = collection.InsertOne(context.TODO(), user)
		if err != nil {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User successfully registered",
		})
	}
}
