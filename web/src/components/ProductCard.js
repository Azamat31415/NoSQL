import React from "react";
import { useNavigate } from "react-router-dom";
import "./ProductCard.css";

const ProductCard = ({ product }) => {
    const navigate = useNavigate();
    const userRole = localStorage.getItem("role");
    const token = localStorage.getItem("token");

    const handleAddToCart = async () => {
        const userID = Number(localStorage.getItem("userID"));  // 🔥 Исправлено
        const token = localStorage.getItem("token");

        console.log("Отправляемый токен:", token);
        console.log("Отправляемый userID:", userID);  // Теперь точно число!

        if (!userID) {
            alert("Please log in to add items to your cart.");
            return;
        }

        if (!token) {
            alert("Authorization token is missing. Please log in again.");
            return;
        }

        if (!product.id) {
            console.error("Ошибка: product.id отсутствует!", product);
            alert("Invalid product ID. Please try again later.");
            return;
        }

        const cartItem = {
            user_id: userID,  // Теперь это число
            product_id: product.id,
            quantity: 1
        };

        console.log("Отправляем в корзину:", cartItem);

        try {
            const response = await fetch("http://localhost:8080/cart", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify(cartItem),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Failed to add item to cart: ${errorText}`);
            }

            alert("Item added to cart!");
        } catch (error) {
            console.error("Error:", error);
            alert("Error adding item to cart");
        }
    };

    const handleDeleteProduct = async () => {
        if (!window.confirm("Are you sure you want to delete this product?")) return;

        if (!product.id) { // Исправлено
            console.error("Ошибка: product.id отсутствует!");
            alert("Invalid product ID. Please try again later.");
            return;
        }

        try {
            const response = await fetch(`http://localhost:8080/products/${product.id}`, { // Исправлено
                method: "DELETE",
                headers: {
                    "Authorization": `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Failed to delete product: ${errorText}`);
            }

            alert("Product deleted successfully!");
            window.location.reload();
        } catch (error) {
            console.error("Error:", error);
            alert("Error deleting product");
        }
    };

    return (
        <div className="product-card">
            <img
                src={product.image || "https://catspaw.ru/wp-content/uploads/2016/06/Ela_Kaimo.jpg"}
                alt={product.name || "Unavailable"}
                className="product-image"
            />
            <h3 className="product-name">{product.name || "Unavailable"}</h3>
            <p className="product-description">{product.description || "No description available"}</p>
            <p className="product-price">Price: {product.price || "Not available"}</p>
            <button className="add-to-cart-button" onClick={handleAddToCart}>
                Add to Cart
            </button>
            {userRole === "admin" && (
                <>
                    <p className="product-id"><strong>ID:</strong> {product.id || "N/A"}</p>
                    <button className="edit-product-button" onClick={() => navigate(`/edit-product/${product.id}`)}>
                        Edit
                    </button>
                    <button className="delete-product-button" onClick={handleDeleteProduct}>
                        Delete
                    </button>
                </>
            )}
        </div>
    );
};

export default ProductCard;
