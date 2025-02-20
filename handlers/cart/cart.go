package cart

import (
	"GoProject/config"
	"GoProject/migrations"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

var cartCollection *mongo.Collection

func InitCartCollection(db *mongo.Database) {
	cartCollection = db.Collection("cart")
}

func AddToCart(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := config.VerifyJWT(w, r, db.Collection("users"))
		if err != nil || user == nil {
			return
		}

		var cartItem struct {
			ProductID string `json:"product_id"`
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

		userID := user.ID
		productID, err := primitive.ObjectIDFromHex(cartItem.ProductID)
		if err != nil {
			http.Error(w, "Invalid ProductID", http.StatusBadRequest)
			return
		}

		cartCollection := db.Collection("cart")
		filter := bson.M{"user_id": userID, "product_id": productID}
		var existingCartItem migrations.CartItem
		err = cartCollection.FindOne(context.TODO(), filter).Decode(&existingCartItem)

		if err == nil {
			// Товар уже в корзине — обновляем количество
			update := bson.M{"$inc": bson.M{"quantity": cartItem.Quantity}}
			_, err := cartCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
				return
			}
			existingCartItem.Quantity += cartItem.Quantity
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingCartItem)
			return
		}

		if err == mongo.ErrNoDocuments {
			newCartItem := migrations.CartItem{
				ID:        primitive.NewObjectID(),
				UserID:    userID,
				ProductID: productID,
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

		objectID, err := primitive.ObjectIDFromHex(cartItemID)
		if err != nil {
			http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": objectID}
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

		update := bson.M{"$set": bson.M{"quantity": newQuantity}}
		_, err = cartCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "Failed to update cart item quantity", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Quantity updated"})
	}
}

func RemoveFromCart(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartItemID := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(cartItemID)
		if err != nil {
			http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": objectID}
		_, err = cartCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			http.Error(w, "Failed to remove item from cart", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Item permanently removed from cart"})
	}
}

func GetCartByUser(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")

		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		pipeline := []bson.M{
			{"$match": bson.M{"user_id": objectID}},
			{"$lookup": bson.M{
				"from":         "products",
				"localField":   "product_id",
				"foreignField": "_id",
				"as":           "product",
			}},
			{"$unwind": "$product"},
		}

		cursor, err := cartCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.TODO())

		var cartItems []bson.M
		if err = cursor.All(context.TODO(), &cartItems); err != nil {
			http.Error(w, "Failed to decode cart items", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cartItems)
	}
}

func RemoveOneItemFromCart(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartItemID := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(cartItemID)
		if err != nil {
			http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
			return
		}

		// Уменьшаем количество на 1
		update := bson.M{"$inc": bson.M{"quantity": -1}}
		result, err := cartCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
		if err != nil || result.ModifiedCount == 0 {
			http.Error(w, "Failed to update quantity", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Item quantity decreased"})
	}
}

func GetCartID(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		productID := chi.URLParam(r, "product_id")

		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		productObjectID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var cartItem migrations.CartItem
		err = cartCollection.FindOne(context.TODO(), bson.M{
			"user_id":    userObjectID,
			"product_id": productObjectID,
		}).Decode(&cartItem)

		if err == mongo.ErrNoDocuments {
			http.Error(w, "Cart item not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cartItem)
	}
}
