package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonalPet struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Species string             `bson:"species" json:"species"`
	Age     int                `bson:"age" json:"age"`
	UserID  primitive.ObjectID `bson:"user_id" json:"user_id"`
}

func MigratePersonalPet(db *mongo.Database) error {
	collection := db.Collection("personal_pets")
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"user_id": 1},
		Options: options.Index().SetUnique(false),
	})
	return err
}
