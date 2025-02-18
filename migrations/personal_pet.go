package migrations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PersonalPet struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Species string             `bson:"species"`
	Age     int                `bson:"age"`
	UserID  string             `bson:"user_id"` // user_id теперь строка, а не uint
}

func MigratePersonalPet() error {
	// В MongoDB миграции не требуются в традиционном понимании.
	// Нужно просто создать коллекцию, если она не существует.
	return nil
}
