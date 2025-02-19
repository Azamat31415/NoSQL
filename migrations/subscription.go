package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Subscription struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	StartDate    time.Time          `bson:"start_date" json:"start_date"`
	RenewalDate  time.Time          `bson:"renewal_date" json:"renewal_date"`
	IntervalDays int                `bson:"interval_days" json:"interval_days"`
	Type         string             `bson:"type" json:"type"`
	Status       string             `bson:"status" json:"status"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func MigrateSubscription(db *mongo.Database) error {
	collection := db.Collection("subscriptions")

	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"user_id": 1},
		Options: nil,
	})
	return err
}

func ExpireSubscriptionsNow(db *mongo.Database) error {
	now := time.Now()
	collection := db.Collection("subscriptions")

	_, err := collection.UpdateMany(
		context.TODO(),
		bson.M{
			"renewal_date": bson.M{"$lt": now},
			"status":       "active",
		},
		bson.M{"$set": bson.M{"status": "expired"}},
	)
	return err
}
