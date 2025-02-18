package order

import (
	"GoProject/migrations"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var validStatuses = []string{
	"pending", "shipped", "delivered", "cancelled", "returned",
}

func UpdateOrderStatus(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "id")

		var request struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		isValidStatus := false
		for _, status := range validStatuses {
			if status == request.Status {
				isValidStatus = true
				break
			}
		}

		if !isValidStatus {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		var order migrations.Order
		err := collection.FindOne(r.Context(), bson.M{"_id": orderID}).Decode(&order)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		order.Status = request.Status
		_, err = collection.UpdateOne(
			r.Context(),
			bson.M{"_id": orderID},
			bson.M{"$set": bson.M{"status": order.Status}},
		)
		if err != nil {
			http.Error(w, "Failed to update order status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
	}
}
