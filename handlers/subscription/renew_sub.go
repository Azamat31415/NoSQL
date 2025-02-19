package subscription

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func RenewSubscription(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptionIDStr := chi.URLParam(r, "id")
		if subscriptionIDStr == "" {
			http.Error(w, "Missing subscription ID", http.StatusBadRequest)
			return
		}

		subscriptionID, err := primitive.ObjectIDFromHex(subscriptionIDStr)
		if err != nil {
			http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
			return
		}

		collection := db.Collection("subscriptions")

		var subscription migrations.Subscription
		err = collection.FindOne(context.TODO(), bson.M{"_id": subscriptionID}).Decode(&subscription)
		if err != nil {
			http.Error(w, "Subscription not found", http.StatusNotFound)
			return
		}

		if subscription.Status != "active" {
			http.Error(w, "Cannot renew a non-active subscription", http.StatusBadRequest)
			return
		}

		now := time.Now()
		if subscription.RenewalDate.Before(now) {
			subscription.RenewalDate = now.AddDate(0, 0, subscription.IntervalDays)
		} else {
			subscription.RenewalDate = subscription.RenewalDate.AddDate(0, 0, subscription.IntervalDays)
		}
		subscription.UpdatedAt = now

		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": subscriptionID},
			bson.M{"$set": bson.M{"renewal_date": subscription.RenewalDate, "updated_at": subscription.UpdatedAt}},
		)

		if err != nil {
			http.Error(w, "Failed to renew subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subscription)
	}
}

func ExpireSubscriptionsNowHandler(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := migrations.ExpireSubscriptionsNow(db)
		if err != nil {
			http.Error(w, "Failed to expire subscriptions", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Expired subscriptions successfully updated"))
	}
}
