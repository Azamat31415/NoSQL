package routes

import (
	"GoProject/handlers/auth"
	"GoProject/handlers/cart"
	"GoProject/handlers/subscription"

	//"GoProject/handlers/order"
	//"GoProject/handlers/personal_pet"
	"GoProject/handlers/product"
	//"GoProject/handlers/subscription"
	//"GoProject/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	//"net/http"
)

func InitializeRoutes(r *chi.Mux, db *mongo.Database) {
	// Apply CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes for products
	r.Post("/products", product.AddProduct(db))
	r.Put("/products/{id}", product.EditProduct(db))
	r.Delete("/products/{id}", product.DeleteProduct(db))
	r.Get("/products", product.FetchAllProducts(db))
	r.Get("/products/{id}", product.FetchProductByID(db))
	//
	//// Routes for authentication
	r.Post("/register", auth.RegisterHandler(db.Collection("users")))
	r.Post("/login", auth.LoginHandler(db.Collection("users")))
	r.Get("/profile", auth.ProfileHandler(db.Collection("users")))
	r.Get("/users", auth.GetUsersHandler(db.Collection("users")))

	//// Routes for orders
	//r.Post("/orders", order.CreateOrder(db)) // Change order handlers to work with MongoDB
	//r.Put("/orders/{id}/status/update", order.UpdateOrderStatus(db))
	//r.Put("/orders/{order_id}/delivery", order.ChooseDeliveryMethod(db))
	//r.Get("/orders", order.GetOrders(db))
	//r.Get("/order-history/{user_id}", order.GetOrderHistory(db))
	//
	//// Protected routes using JWT middleware
	//r.Group(func(protected chi.Router) {
	//	protected.Use(middleware.JWTMiddleware)
	//
	//	protected.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("This is a protected route"))
	//	})
	//
	//	protected.Post("/logout", auth.LogoutHandler())
	//})
	//r.Get("/users/{id}/address", auth.GetUserAddress(db))
	//
	//// Routes for pets
	//r.Post("/pets", personal_pet.AddUserPet(db)) // Same for pets CRUD operations
	//r.Put("/pets/{id}", personal_pet.EditUserPet(db))
	//r.Delete("/pets/{id}", personal_pet.DeleteUserPet(db))
	//r.Get("/pets/{id}", personal_pet.FetchUserPets(db))
	//r.Get("/users/{userID}/pets", personal_pet.FetchUserPetByID(db))
	//
	//// Routes for subscriptions
	r.Post("/subscriptions", subscription.CreateSubscription(db))
	r.Delete("/subscriptions/{id}", subscription.DeleteSubscription(db))
	r.Put("/subscriptions/{id}/renew", subscription.RenewSubscription(db))
	r.Post("/subpayment", subscription.HandleSubscriptionPayment(db))
	r.Get("/subscriptions/{user_id}", subscription.GetUserSubscription(db))
	r.Put("/subscriptions/expire", subscription.ExpireSubscriptionsNowHandler(db))
	//
	//// Routes for cart
	r.Post("/cart", cart.AddToCart(db)) // Добавление товара в корзину
	r.Delete("/cart/{id}", cart.RemoveFromCart(db))
	r.Put("/cart/update/{id}/{quantity}", cart.UpdateCartItemQuantity(db))
	r.Delete("/cart/{id}/byone", cart.RemoveOneItemFromCart(db))
	r.Get("/cart/user/{user_id}/products", cart.GetCartByUser(db))
	r.Get("/cart/{user_id}/{product_id}", cart.GetCartID(db))
}
