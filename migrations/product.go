package migrations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Product модель для MongoDB
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	Category    string             `bson:"category"`
	Subcategory string             `bson:"subcategory"`
	Type        string             `bson:"type"`
	CreatedAt   time.Time          `bson:"created_at"`
}
