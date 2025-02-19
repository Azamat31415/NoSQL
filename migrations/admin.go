package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssignAdminRole проверяет и обновляет роль пользователя
func AssignAdminRole(db *mongo.Database) error {
	collection := db.Collection("users")
	email := "azabraza061005@gmail.com"

	ctx := context.TODO()

	// Проверяем, есть ли пользователь
	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// Создаём нового пользователя с ролью "admin"
		adminUser := User{
			ID:        primitive.NewObjectID(),
			Email:     email,
			Password:  "aza061005",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
		}
		if err := adminUser.HashPassword(); err != nil {
			return fmt.Errorf("ошибка хеширования пароля: %v", err)
		}

		_, err := collection.InsertOne(ctx, adminUser)
		if err != nil {
			return fmt.Errorf("ошибка создания администратора: %v", err)
		}
		fmt.Println("Администратор создан.")
		return nil
	} else if err != nil {
		return fmt.Errorf("ошибка при поиске пользователя: %v", err)
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления роли пользователя: %v", err)
	}

	fmt.Println("Роль пользователя обновлена до admin.")
	return nil
}
