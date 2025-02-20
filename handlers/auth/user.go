package auth

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func GetUserAddress(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user migrations.User
		err = collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"address": user.Address,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// GetUsersHandler retrieves all users from the MongoDB collection
func GetUsersHandler(usersCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Query all users from MongoDB
		cursor, err := usersCollection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var users []bson.M
		if err := cursor.All(ctx, &users); err != nil {
			http.Error(w, "Error processing user data", http.StatusInternalServerError)
			return
		}

		// Send JSON response with users
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
