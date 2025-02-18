package cart

import (
	"GoProject/configs"
	"GoProject/migrations"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

var cartCollection *mongo.Collection

func InitCartCollection(db *mongo.Database) {
	cartCollection = db.Collection("cart_items")
}

func AddToCart(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := config.VerifyJWT(w, r, db)
		if err != nil || user == nil {
			return // Ошибка уже обработана в VerifyJWT
		}

		var cartItem struct {
			ProductID string `json:"product_id"` // Преобразуем в строку для MongoDB
			Quantity  int    `json:"quantity"`
		}

		if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if cartItem.ProductID == "" || cartItem.Quantity <= 0 {
			http.Error(w, "ProductID and valid Quantity are required", http.StatusBadRequest)
			return
		}

		filter := bson.M{"user_id": user.ID, "product_id": cartItem.ProductID}
		var existingCartItem migrations.CartItem
		err = cartCollection.FindOne(context.TODO(), filter).Decode(&existingCartItem)

		if err == nil {
			existingCartItem.Quantity += cartItem.Quantity
			_, err := cartCollection.ReplaceOne(context.TODO(), filter, existingCartItem)
			if err != nil {
				http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingCartItem)
			return
		}

		if err == mongo.ErrNoDocuments {
			newCartItem := migrations.CartItem{
				UserID:    user.ID,
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
			}

			_, err := cartCollection.InsertOne(context.TODO(), newCartItem)
			if err != nil {
				http.Error(w, "Failed to add item to cart", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCartItem)
			return
		}

		http.Error(w, "Database error", http.StatusInternalServerError)
	}
}

func UpdateCartItemQuantity(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartItemID := chi.URLParam(r, "id")
		quantityStr := chi.URLParam(r, "quantity")

		newQuantity, err := strconv.Atoi(quantityStr)
		if err != nil || newQuantity < 0 {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": cartItemID}
		var existingCartItem migrations.CartItem
		err = cartCollection.FindOne(context.TODO(), filter).Decode(&existingCartItem)
		if err != nil {
			http.Error(w, "Item not found in cart", http.StatusNotFound)
			return
		}

		if newQuantity == 0 {
			_, err := cartCollection.DeleteOne(context.TODO(), filter)
			if err != nil {
				http.Error(w, "Failed to remove item from cart", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart"})
			return
		}

		existingCartItem.Quantity = newQuantity
		_, err = cartCollection.ReplaceOne(context.TODO(), filter, existingCartItem)
		if err != nil {
			http.Error(w, "Failed to update cart item quantity", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingCartItem)
	}
}

func RemoveFromCart(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartItemID := chi.URLParam(r, "id")
		if cartItemID == "" {
			http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": cartItemID}
		_, err := cartCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			http.Error(w, "Failed to remove item from cart", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Item permanently removed from cart"})
	}
}
