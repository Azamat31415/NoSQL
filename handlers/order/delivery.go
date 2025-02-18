package order

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type ChooseDeliveryRequest struct {
	DeliveryMethod string  `json:"delivery_method" binding:"required"`
	PickupPointID  *string `json:"pickup_point_id,omitempty"`
	Address        *string `json:"address,omitempty"`
}

func ChooseDeliveryMethod(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChooseDeliveryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %s", err.Error()), http.StatusBadRequest)
			return
		}

		orderID := chi.URLParam(r, "order_id")
		var order migrations.Order

		// Находим заказ по ID
		err := collection.FindOne(context.TODO(), bson.M{"_id": orderID}).Decode(&order)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		order.DeliveryMethod = req.DeliveryMethod
		if req.DeliveryMethod == "pickup" {
			if req.PickupPointID != nil {
				order.PickupPointID = req.PickupPointID
			} else {
				http.Error(w, "Pickup point ID is required for pickup", http.StatusBadRequest)
				return
			}
		} else if req.DeliveryMethod == "delivery" {
			if req.Address == nil {
				http.Error(w, "Address is required for delivery", http.StatusBadRequest)
				return
			}
			order.Address = req.Address
		}

		// Обновление заказа в MongoDB
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": orderID},
			bson.M{"$set": bson.M{
				"delivery_method": order.DeliveryMethod,
				"pickup_point_id": order.PickupPointID,
				"address":         order.Address,
			}},
		)
		if err != nil {
			http.Error(w, "Failed to update delivery method", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"message": "Delivery method updated successfully",
			"order":   order,
		}
		json.NewEncoder(w).Encode(response)
	}
}
