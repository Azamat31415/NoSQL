package personal_pet

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
)

// Получение коллекции MongoDB
func getPersonalPetCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("personal_pets")
}

func AddUserPet(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pet migrations.PersonalPet

		// Декодируем тело запроса в структуру pet
		if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Проверка на обязательные данные
		if pet.UserID == "" {
			http.Error(w, "UserID is required", http.StatusBadRequest)
			return
		}

		fmt.Printf("Saving pet: %+v\n", pet)

		// Вставляем новый питомец в коллекцию
		collection := getPersonalPetCollection(db)
		pet.ID = primitive.NewObjectID() // Генерация нового ObjectID

		_, err := collection.InsertOne(context.TODO(), pet)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add pet: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pet)
	}
}

func EditUserPet(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var pet migrations.PersonalPet

		// Преобразуем ID из строки в ObjectID
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid pet ID", http.StatusBadRequest)
			return
		}

		// Ищем питомца по ID
		collection := getPersonalPetCollection(db)
		filter := bson.M{"_id": objectID}
		err = collection.FindOne(context.TODO(), filter).Decode(&pet)
		if err != nil {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}

		// Декодируем новые данные
		if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Обновляем питомца в базе данных
		_, err = collection.ReplaceOne(context.TODO(), filter, pet)
		if err != nil {
			http.Error(w, "Failed to update pet", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pet)
	}
}

func DeleteUserPet(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var pet migrations.PersonalPet

		// Преобразуем ID из строки в ObjectID
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid pet ID", http.StatusBadRequest)
			return
		}

		// Удаляем питомца из коллекции
		collection := getPersonalPetCollection(db)
		filter := bson.M{"_id": objectID}
		_, err = collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			http.Error(w, "Failed to delete pet", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pet deleted successfully"))
	}
}

func FetchUserPets(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var pet migrations.PersonalPet

		// Преобразуем ID из строки в ObjectID
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid pet ID", http.StatusBadRequest)
			return
		}

		// Ищем питомца по ID
		collection := getPersonalPetCollection(db)
		filter := bson.M{"_id": objectID}
		err = collection.FindOne(context.TODO(), filter).Decode(&pet)
		if err != nil {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pet)
	}
}

func FetchUserPetByID(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		var pets []migrations.PersonalPet

		// Ищем всех питомцев для данного userID
		collection := getPersonalPetCollection(db)
		filter := bson.M{"user_id": userID}
		cursor, err := collection.Find(context.TODO(), filter)
		if err != nil {
			http.Error(w, "Pets not found", http.StatusNotFound)
			return
		}

		// Чтение всех найденных питомцев
		if err := cursor.All(context.TODO(), &pets); err != nil {
			http.Error(w, "Failed to decode pets", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pets)
	}
}
