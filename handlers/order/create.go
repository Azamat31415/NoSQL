package order

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type CreateOrderRequest struct {
	UserID         string             `json:"user_id"`
	DeliveryMethod string             `json:"delivery_method"`
	PickupPointID  *string            `json:"pickup_point_id"`
	Address        *string            `json:"address"`
	TotalPrice     float64            `json:"total_price"`
	OrderItems     []OrderItemRequest `json:"order_items"`
}

type OrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func CreateOrder(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Создание нового заказа
		order := migrations.Order{
			UserID:         request.UserID,
			DeliveryMethod: request.DeliveryMethod,
			PickupPointID:  request.PickupPointID,
			Address:        request.Address,
			Status:         "pending",
			TotalPrice:     request.TotalPrice,
			OrderItems:     []migrations.OrderItem{}, // временно пустой список
			CreatedAt:      time.Now(),
		}

		// Вставка заказа в MongoDB
		result, err := collection.InsertOne(context.TODO(), order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Проверка типа InsertedID, чтобы корректно извлечь ObjectID
		var orderID primitive.ObjectID
		switch v := result.InsertedID.(type) {
		case primitive.ObjectID:
			orderID = v
		default:
			http.Error(w, "Invalid InsertedID type", http.StatusInternalServerError)
			return
		}

		// Добавление элементов заказа
		for _, item := range request.OrderItems {
			orderItem := migrations.OrderItem{
				OrderID:   orderID, // Используем корректный ObjectID
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
				CreatedAt: time.Now(),
			}

			// Коллекция для сохранения элементов заказа
			orderItemsCollection := collection.Database().Collection("order_items")
			_, err := orderItemsCollection.InsertOne(context.TODO(), orderItem)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Обновляем основной заказ, добавляя к нему элементы заказа (опционально)
		_, err = collection.UpdateOne(
			context.TODO(),
			primitive.M{"_id": orderID},
			primitive.M{"$set": primitive.M{"order_items": request.OrderItems}},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправка успешного ответа с добавленным заказом
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}
