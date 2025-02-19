package main

import (
	"GoProject/config"
	"GoProject/handlers/cart" // Добавляем импорт cart
	"GoProject/migrations"
	"GoProject/routes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	// Подключение к MongoDB
	db, err := config.GetMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Запуск миграций (если нужны)
	if err := migrations.MigrateAll(db); err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	// Назначение админ-ролей
	if err := migrations.AssignAdminRole(db); err != nil {
		log.Printf("Failed to assign admin role: %v", err)
	} else {
		fmt.Println("Admin role assignment completed.")
	}

	// Инициализация коллекции корзины
	cart.InitCartCollection(db)

	// Инициализация маршрутов
	r := chi.NewRouter()

	// CORS настройки
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Используем CORS middleware
	r.Use(c.Handler)

	// Инициализация маршрутов
	routes.InitializeRoutes(r, db)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
