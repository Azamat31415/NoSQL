package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionPayment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SubscriptionID primitive.ObjectID `bson:"subscription_id" json:"subscription_id"`
	Amount         float64            `bson:"amount" json:"amount"`
	PaymentDate    time.Time          `bson:"payment_date" json:"payment_date"`
	Status         string             `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

func MigrateSubscriptionPayment(db *mongo.Database) error {
	collection := db.Collection("subscription_payments")
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"subscription_id": 1},
		Options: nil,
	})
	return err
}
