import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./CartPage.css";

const CartPage = () => {
    const [cart, setCart] = useState([]);
    const [selectedItems, setSelectedItems] = useState([]);
    const navigate = useNavigate();

    useEffect(() => {
        fetchCartItems();
    }, []);

    const fetchCartItems = async () => {
        const token = localStorage.getItem("token");
        const userID = localStorage.getItem("userID");

        if (!token || !userID) {
            console.error("Missing token or userID");
            return;
        }

        try {
            const response = await fetch(`http://localhost:8080/cart/user/${userID}/products`, {
                method: "GET",
                headers: { Authorization: `Bearer ${token}` },
            });

            if (!response.ok) throw new Error("Failed to fetch cart items");

            const data = await response.json();
            console.log("Cart items received from API:", data);
            setCart(data || []);
        } catch (error) {
            console.error("Error fetching cart:", error);
            setCart([]);
        }
    };

    const updateQuantity = async (id, quantity) => {
        if (quantity < 1) return;

        try {
            const token = localStorage.getItem("token");
            const response = await fetch(`http://localhost:8080/cart/${id}/quantity/${quantity}`, {
                method: "PUT",
                headers: { Authorization: `Bearer ${token}` },
            });

            if (!response.ok) throw new Error("Failed to update quantity");

            setCart(cart.map((item) => (item._id === id ? { ...item, quantity } : item)));
        } catch (error) {
            console.error("Error updating quantity:", error);
        }
    };

    const removeFromCart = async (id) => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch(`http://localhost:8080/cart/${id}`, {
                method: "DELETE",
                headers: { Authorization: `Bearer ${token}` },
            });

            if (!response.ok) throw new Error("Failed to remove item from cart");

            setCart(cart.filter((item) => item._id !== id));
        } catch (error) {
            console.error("Error removing item from cart:", error);
        }
    };

    const toggleSelection = (id) => {
        setSelectedItems((prevSelected) =>
            prevSelected.includes(id)
                ? prevSelected.filter((item) => item !== id)
                : [...prevSelected, id]
        );
    };

    return (
        <div className="cart-page">
            <h2>Your Cart</h2>
            {cart.length > 0 ? (
                <div className="cart-list">
                    {cart.map((item) => (
                        <div key={item._id} className="cart-item">
                            <input
                                type="checkbox"
                                checked={selectedItems.includes(item._id)}
                                onChange={() => toggleSelection(item._id)}
                                className="cart-checkbox"
                            />

                            <img
                                src={item.image ? item.image : "placeholder.jpg"}
                                alt={item.name ? item.name : "No image"}
                                className="cart-item-image"
                            />

                            <div className="cart-item-info">
                                <div className="cart-item-details">
                                    <h3>{item.name ? item.name : "No name"}</h3>
                                    <p>Price: ${item.price ? item.price.toFixed(2) : "0.00"}</p>
                                    <p>Total: ${(item.price && item.quantity) ? (item.price * item.quantity).toFixed(2) : "0.00"}</p>

                                    <div className="quantity-control">
                                        <button onClick={() => updateQuantity(item._id, item.quantity - 1)}>-</button>
                                        <input
                                            type="number"
                                            value={item.quantity ?? 1}
                                            min="1"
                                            onChange={(e) => updateQuantity(item._id, parseInt(e.target.value))}
                                        />
                                        <button onClick={() => updateQuantity(item._id, item.quantity + 1)}>+</button>
                                    </div>
                                </div>

                                <button className="remove-button" onClick={() => removeFromCart(item._id)}>Remove</button>
                            </div>
                        </div>
                    ))}
                    <button
                        className="proceed-button"
                        onClick={() => {
                            localStorage.setItem("selectedItems", JSON.stringify(selectedItems));
                            navigate("/payment");
                        }}
                        disabled={selectedItems.length === 0} // Блокируем кнопку, если ничего не выбрано
                    >
                        Proceed to Payment
                    </button>
                </div>
            ) : (
                <p>Your cart is empty.</p>
            )}
        </div>
    );
};

export default CartPage;
