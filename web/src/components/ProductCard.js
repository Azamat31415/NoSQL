import React from "react";
import { useNavigate } from "react-router-dom";
import "./ProductCard.css";

const ProductCard = ({ product }) => {
    const navigate = useNavigate();
    const userRole = localStorage.getItem("role");
    const token = localStorage.getItem("token");

    const handleAddToCart = async () => {
        const userID = localStorage.getItem("userID");

        if (!userID || isNaN(parseInt(userID))) {
            alert("Please log in to add items to your cart.");
            return;
        }

        if (!token) {
            alert("Authorization token is missing. Please log in again.");
            return;
        }

        const cartItem = {
            user_id: parseInt(userID),
            product_id: product.ID,
            quantity: 1
        };

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
                throw new Error("Failed to add item to cart");
            }

            alert("Item added to cart!");
        } catch (error) {
            console.error("Error:", error);
            alert("Error adding item to cart");
        }
    };

    const handleDeleteProduct = async () => {
        if (!window.confirm("Are you sure you want to delete this product?")) return;

        try {
            const response = await fetch(`http://localhost:8080/products/${product.ID}`, {
                method: "DELETE",
                headers: {
                    "Authorization": `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                throw new Error("Failed to delete product");
            }

            alert("Product deleted successfully!");
            window.location.reload(); // Refresh the page to update the product list
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
                    <p className="product-id"><strong>ID:</strong> {product.ID}</p>
                    <button className="edit-product-button" onClick={() => navigate(`/edit-product/${product.ID}`)}>
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