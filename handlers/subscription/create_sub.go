package subscription

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateSubscriptionRequest struct {
	UserID       string `json:"user_id"`
	IntervalDays int    `json:"interval_days"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

func CreateSubscription(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		fmt.Println("Received UserID:", req.UserID) // Логирование

		userID, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		subscription := bson.M{
			"_id":           primitive.NewObjectID(),
			"user_id":       userID,
			"start_date":    time.Now(),
			"renewal_date":  time.Now().AddDate(0, 0, req.IntervalDays),
			"interval_days": req.IntervalDays,
			"type":          req.Type,
			"status":        req.Status,
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		}

		collection := db.Collection("subscriptions")
		_, err = collection.InsertOne(context.TODO(), subscription)
		if err != nil {
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(subscription)
	}
}

func DeleteSubscription(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptionID := chi.URLParam(r, "id")
		id, err := primitive.ObjectIDFromHex(subscriptionID)
		if err != nil {
			http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
			return
		}

		collection := db.Collection("subscriptions")
		_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleSubscriptionPayment(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SubscriptionID string  `json:"subscription_id"`
			Amount         float64 `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		subID, err := primitive.ObjectIDFromHex(req.SubscriptionID)
		if err != nil {
			http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
			return
		}

		payment := migrations.SubscriptionPayment{
			ID:             primitive.NewObjectID(),
			SubscriptionID: subID,
			Amount:         req.Amount,
			PaymentDate:    time.Now(),
			Status:         "Paid",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		collection := db.Collection("subscription_payments")
		_, err = collection.InsertOne(context.TODO(), payment)
		if err != nil {
			http.Error(w, "Failed to process payment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(payment)
	}
}

func GetUserSubscription(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		collection := db.Collection("subscriptions")

		var subscription migrations.Subscription
		err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&subscription)
		if err != nil {
			fmt.Println("Subscription not found for user:", userID) // Добавь этот лог
			http.Error(w, "No active subscription found.", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(subscription)
	}
}
