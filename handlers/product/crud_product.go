package product

import (
	"GoProject/migrations"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// ProductResponse — формат для JSON-ответа (чтобы ObjectID был строкой)
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
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Генерируем новый ObjectID
		product.ID = primitive.NewObjectID()

		collection := db.Collection("products")
		_, err := collection.InsertOne(context.Background(), product)
		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		// Форматируем ответ
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

// EditProduct — обновляет продукт
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

		// Убираем пустые поля из обновления
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

// DeleteProduct удаляет продукт по ID
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

// FetchAllProducts — получает список всех продуктов
func FetchAllProducts(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := db.Collection("products")

		cursor, err := collection.Find(context.Background(), bson.M{})
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

		// Форматируем JSON-ответ (конвертируем ObjectID в строку)
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

// FetchProductByID — получает продукт по ID
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

		// Форматируем ответ
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
