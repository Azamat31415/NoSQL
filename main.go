package main

import (
	"GoProject/config"
	"GoProject/migrations"
	"GoProject/routes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

// Background goroutine to check expired subscriptions
func checkExpiredSubscriptions(db *mongo.Database) {
	for {
		time.Sleep(24 * time.Hour)

		now := time.Now()

		// MongoDB query to find expired subscriptions
		filter := bson.M{"renewal_date": bson.M{"$lt": now}, "status": "active"}
		update := bson.M{"$set": bson.M{"status": "expired"}}

		// Update all matching subscriptions
		_, err := db.Collection("subscriptions").UpdateMany(context.TODO(), filter, update)
		if err != nil {
			log.Println("Error updating subscriptions:", err)
		} else {
			log.Println("Expired subscriptions updated.")
		}
	}
}

func main() {
	// Connect to MongoDB
	db, err := config.GetMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Run migrations (if needed for MongoDB)
	if err := migrations.MigrateAll(db); err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	// Assign admin role (if required)
	if err := migrations.AssignAdminRole(db); err != nil {
		log.Printf("Failed to assign admin role: %v", err)
	} else {
		fmt.Println("Admin role assignment completed.")
	}

	// Start background subscription check
	go checkExpiredSubscriptions(db)

	// Initialize chi router
	r := chi.NewRouter()

	// CORS settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Use CORS middleware
	r.Use(c.Handler)

	// Initialize routes with MongoDB connection
	routes.InitializeRoutes(r, db)

	// Start the server
	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
