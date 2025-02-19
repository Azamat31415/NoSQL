package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Order struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID  `bson:"user_id" json:"user_id"`
	DeliveryMethod string              `bson:"delivery_method" json:"delivery_method"`
	PickupPointID  *primitive.ObjectID `bson:"pickup_point_id,omitempty" json:"pickup_point_id,omitempty"`
	Address        *string             `bson:"address,omitempty" json:"address,omitempty"`
	Status         string              `bson:"status" json:"status"`
	TotalPrice     float64             `bson:"total_price" json:"total_price"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at" json:"updated_at"`
	OrderItems     []OrderItem         `bson:"order_items" json:"order_items"`
}

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func MigrateOrder(db *mongo.Database) error {
	collection := db.Collection("orders")
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"user_id": 1},
		Options: options.Index().SetUnique(false),
	})
	return err
}
