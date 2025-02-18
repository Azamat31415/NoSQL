package subscription

import (
	"GoProject/migrations"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"strconv"
)

func RenewSubscription(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptionIDStr := chi.URLParam(r, "id")
		if subscriptionIDStr == "" {
			http.Error(w, "Missing subscription ID", http.StatusBadRequest)
			return
		}

		subscriptionID, err := strconv.ParseUint(subscriptionIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
			return
		}

		var subscription migrations.Subscription
		if err := db.First(&subscription, subscriptionID).Error; err != nil {
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

		if err := db.Model(&subscription).Updates(map[string]interface{}{
			"renewal_date": subscription.RenewalDate,
			"updated_at":   subscription.UpdatedAt,
		}).Error; err != nil {
			http.Error(w, "Failed to renew subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subscription)
	}
}

func ExpireSubscriptionsNowHandler(db *gorm.DB) http.HandlerFunc {
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
