package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssignAdminRole checks and updates the user's role
func AssignAdminRole(db *mongo.Database) error {
	collection := db.Collection("users")
	email := "azabraza061005@gmail.com"

	ctx := context.TODO()

	// Check if the user exists
	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// Create a new user with the "admin" role
		adminUser := User{
			ID:        primitive.NewObjectID(),
			Email:     email,
			Password:  "aza061005",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
		}
		if err := adminUser.HashPassword(); err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}

		_, err := collection.InsertOne(ctx, adminUser)
		if err != nil {
			return fmt.Errorf("error creating admin: %v", err)
		}
		fmt.Println("Admin user created.")
		return nil
	} else if err != nil {
		return fmt.Errorf("error finding user: %v", err)
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	if err != nil {
		return fmt.Errorf("error updating user role: %v", err)
	}

	fmt.Println("User role updated to admin.")
	return nil
}
