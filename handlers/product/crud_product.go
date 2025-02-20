package product

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory"`
	Type        string  `json:"type"`
}

func AddProduct(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product migrations.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			log.Println("Ошибка декодирования JSON:", err)
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		product.ID = primitive.NewObjectID()

		collection := db.Collection("products")
		_, err := collection.InsertOne(context.Background(), product)
		if err != nil {
			log.Println("Ошибка при добавлении продукта в базу:", err)
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		response := ProductResponse{
			ID:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			Subcategory: product.Subcategory,
			Type:        product.Type,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func EditProduct(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var updatedProduct map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		updateData := bson.M{}
		for key, value := range updatedProduct {
			if value != "" && value != 0 {
				updateData[key] = value
			}
		}

		if len(updateData) == 0 {
			http.Error(w, "No valid fields to update", http.StatusBadRequest)
			return
		}

		collection := db.Collection("products")
		res, err := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": objectID},
			bson.M{"$set": updateData},
		)

		if err != nil || res.MatchedCount == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
	}
}

func DeleteProduct(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		collection := db.Collection("products")
		res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
		if err != nil || res.DeletedCount == 0 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
	}
}

func FetchAllProducts(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := db.Collection("products")

		category := r.URL.Query().Get("category")
		subcategory := r.URL.Query().Get("subcategory")
		productType := r.URL.Query().Get("type")

		filter := bson.M{}
		if category != "" {
			filter["category"] = category
		}
		if subcategory != "" {
			filter["subcategory"] = subcategory
		}
		if productType != "" {
			filter["type"] = productType
		}

		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			http.Error(w, "Error fetching products", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var products []migrations.Product
		if err = cursor.All(context.Background(), &products); err != nil {
			http.Error(w, "Error decoding products", http.StatusInternalServerError)
			return
		}

		var formattedProducts []ProductResponse
		for _, p := range products {
			formattedProducts = append(formattedProducts, ProductResponse{
				ID:          p.ID.Hex(),
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Stock:       p.Stock,
				Category:    p.Category,
				Subcategory: p.Subcategory,
				Type:        p.Type,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(formattedProducts)
	}
}

func FetchProductByID(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		collection := db.Collection("products")

		var product migrations.Product
		err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&product)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		response := ProductResponse{
			ID:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			Subcategory: product.Subcategory,
			Type:        product.Type,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
