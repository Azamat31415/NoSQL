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

            setCart(cart.map((item) => (item.id === id ? { ...item, quantity } : item)));
        } catch (error) {
            console.error("Error updating quantity:", error);
        }
    };

    const toggleSelection = (cartId) => {
        setSelectedItems((prevSelected) =>
            prevSelected.includes(cartId)
                ? prevSelected.filter((id) => id !== cartId)
                : [...prevSelected, cartId]
        );
    };

    const removeFromCart = async (productID) => {
        try {
            const token = localStorage.getItem("token");

            if (!token) {
                console.error("Missing token");
                return;
            }
            const cartId = await getCartItemID(productID);

            if (!cartId) {
                console.error("Cart item not found");
                return;
            }

            console.log(`Deleting cart item with Cart ID: ${cartId}`);

            const response = await fetch(`http://localhost:8080/cart/${cartId}`, {
                method: "DELETE",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (!response.ok) throw new Error("Failed to remove item");

            setCart(cart.filter((item) => item.cart_id !== cartId));
        } catch (error) {
            console.error("Error removing item:", error);
        }
    };

    const getCartItemID = async (productID) => {
        try {
            const token = localStorage.getItem("token");
            const userID = localStorage.getItem("userID");

            const response = await fetch(`http://localhost:8080/cart/${userID}/${productID}`, {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (!response.ok) throw new Error("Failed to fetch cart item ID");

            const data = await response.json();
            console.log("Cart item ID:", data);
            return data.cart_id;
        } catch (error) {
            console.error("Error fetching cart item ID:", error);
        }
    };

    const proceedToCheckout = () => {
        localStorage.setItem("selectedItems", JSON.stringify(selectedItems));
        navigate("/payment");
    };

    const totalPrice = cart.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const roundedTotal = Math.ceil(totalPrice * 100) / 100;

    return (
        <div className="cart-page">
            <h2>Your Cart</h2>
            {cart.length > 0 ? (
                <>
                    <div className="cart-list">
                        {cart.map((item) => (
                            <div key={item.id} className="cart-item">
                                <input
                                    type="checkbox"
                                    checked={selectedItems.includes(item.id)}
                                    onChange={() => toggleSelection(item.id)}
                                    className="cart-checkbox"
                                />
                                <div className="cart-item-details">
                                    <h3>{item.name}</h3>
                                    <p>Price: ${item.price}</p>
                                    <p>Total: ${Math.ceil(item.price * item.quantity * 100) / 100}</p>
                                    <div className="quantity-control">
                                        <button onClick={() => updateQuantity(item.id, item.quantity - 1)}>-</button>
                                        <input
                                            type="number"
                                            value={item.quantity}
                                            min="1"
                                            onChange={(e) => updateQuantity(item.id, parseInt(e.target.value))}
                                        />
                                        <button onClick={() => updateQuantity(item.id, item.quantity + 1)}>+</button>
                                    </div>
                                </div>
                                <button className="remove-button" onClick={() => removeFromCart(item.id)}>
                                    Remove
                                </button>
                            </div>
                        ))}
                    </div>
                    <h3>Total: ${roundedTotal}</h3>
                    <button
                        className="checkout-button"
                        disabled={selectedItems.length === 0}
                        onClick={proceedToCheckout}
                    >
                        Proceed to Payment ({selectedItems.length})
                    </button>
                </>
            ) : (
                <p>Your cart is empty.</p>
            )}
        </div>
    );
};

export default CartPage;
