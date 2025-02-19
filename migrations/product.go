package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product model for MongoDB
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	Category    string             `bson:"category" json:"category"`
	Subcategory string             `bson:"subcategory" json:"subcategory"`
	Type        string             `bson:"type" json:"type"`
}

// MigrateProduct ensures indexes are created
func MigrateProduct(db *mongo.Database) error {
	collection := db.Collection("products")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Printf("Failed to create unique index on 'name': %v\n", err)
	}
	return err
}
