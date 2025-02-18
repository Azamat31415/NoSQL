package migrations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItem представляет элемент корзины для MongoDB
type CartItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID
	UserID    string             `bson:"user_id"`       // ID пользователя
	ProductID string             `bson:"product_id"`    // ID продукта
	Quantity  int                `bson:"quantity"`      // Количество товара
}
