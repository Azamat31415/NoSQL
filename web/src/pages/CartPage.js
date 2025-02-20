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

    return (
        <div className="cart-page">
            <h2>Your Cart</h2>
            {cart.length > 0 ? (
                <>
                    <div className="cart-list">
                        {cart.map((item) => (
                            <div key={item._id} className="cart-item">
                                <input
                                    type="checkbox"
                                    checked={selectedItems.includes(item._id)}
                                    onChange={() => toggleSelection(item._id)}
                                    className="cart-checkbox"
                                />
                                <div className="cart-item-details">
                                    <h3>{item.name}</h3>
                                    <p>Price: ${item.price}</p>
                                    <p>Total: ${Math.ceil(item.price * item.quantity * 100) / 100}</p>
                                    <div className="quantity-control">
                                        <button onClick={() => updateQuantity(item._id, item.quantity - 1)}>-</button>
                                        <input
                                            type="number"
                                            value={item.quantity}
                                            min="1"
                                            onChange={(e) => updateQuantity(item._id, parseInt(e.target.value))}
                                        />
                                        <button onClick={() => updateQuantity(item._id, item.quantity + 1)}>+</button>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </>
            ) : (
                <p>Your cart is empty.</p>
            )}
        </div>
    );
};

export default CartPage;
