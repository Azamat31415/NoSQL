package product

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// AddProduct handles a query for adding a product
func AddProduct(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
			return
		}

		var product migrations.Product
		fmt.Println("Received request to add product")

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		product.ID = primitive.NewObjectID()
		product.CreatedAt = time.Now()
		_, err := collection.InsertOne(context.TODO(), product)
		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		fmt.Println("Product added:", product)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

// DeleteProduct handles deleting a product
func DeleteProduct(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
		if err != nil || res.DeletedCount == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	}
}

// EditProduct handles updating a product
func EditProduct(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var updatedProduct migrations.Product
		if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		update := bson.M{"$set": updatedProduct}
		res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
		if err != nil || res.MatchedCount == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedProduct)
	}
}

// FetchProductByID retrieves a product by ID
func FetchProductByID(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var product migrations.Product
		err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&product)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// FetchAllProducts retrieves all products
func FetchAllProducts(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []migrations.Product
		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			http.Error(w, "Error fetching products", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.TODO())

		if err = cursor.All(context.TODO(), &products); err != nil {
			http.Error(w, "Error decoding products", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}
