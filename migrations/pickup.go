package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PickupPoint struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Address      string             `bson:"address" json:"address"`
	City         string             `bson:"city" json:"city"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Phone        string             `bson:"phone,omitempty" json:"phone,omitempty"`
	WorkingHours string             `bson:"working_hours,omitempty" json:"working_hours,omitempty"`
}

func MigratePickupPoint(db *mongo.Database) error {
	collection := db.Collection("pickup_points")
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"city": 1},
		Options: options.Index().SetUnique(false),
	})
	return err
}
