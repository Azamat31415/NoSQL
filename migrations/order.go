package migrations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Order represents the structure of an order
type Order struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`             // MongoDB ID (если используется строка ObjectID)
	UserID         string             `bson:"user_id"`                   // ID пользователя
	DeliveryMethod string             `bson:"delivery_method"`           // Метод доставки
	PickupPointID  *string            `bson:"pickup_point_id,omitempty"` // ID точки самовывоза (если есть)
	Address        *string            `bson:"address,omitempty"`         // Адрес (если есть)
	Status         string             `bson:"status"`                    // Статус заказа
	TotalPrice     float64            `bson:"total_price"`               // Общая цена
	CreatedAt      time.Time          `bson:"created_at"`                // Дата создания
	UpdatedAt      time.Time          `bson:"updated_at"`                // Дата обновления
	OrderItems     []OrderItem        `bson:"order_items"`               // Список товаров в заказе
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ID
	OrderID   primitive.ObjectID `bson:"order_id"`      // ID заказа
	ProductID string             `bson:"product_id"`    // ID продукта
	Quantity  int                `bson:"quantity"`      // Количество
	Price     float64            `bson:"price"`         // Цена за единицу
	CreatedAt time.Time          `bson:"created_at"`    // Дата добавления товара
}
