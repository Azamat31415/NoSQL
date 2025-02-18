package main

import (
	"GoProject/configs"
	"GoProject/migrations"
	"GoProject/routes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// Background goroutine to check expired subscriptions
func checkExpiredSubscriptions(db *gorm.DB) {
	for {
		time.Sleep(24 * time.Hour)

		now := time.Now()
		result := db.Model(&migrations.Subscription{}).
			Where("renewal_date < ? AND status = ?", now, "active").
			Update("status", "expired")

		if result.Error != nil {
			log.Println("Error updating subscriptions:", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Updated %d subscriptions: status changed to 'expired'", result.RowsAffected)
		}
	}
}

func main() {
	// Database connection settings
	dsn := config.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Run migrations
	if err := migrations.MigrateAll(db); err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	// Assign admin role
	if err := migrations.AssignAdminRole(db); err != nil {
		log.Printf("Failed to assign admin role: %v", err)
	} else {
		fmt.Println("Admin role assignment completed.")
	}

	// Start background subscription check
	go checkExpiredSubscriptions(db)

	// Initialize chi router
	r := chi.NewRouter()

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(c.Handler)

	// Initialize routes
	routes.InitializeRoutes(r, db)

	// Start server
	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
