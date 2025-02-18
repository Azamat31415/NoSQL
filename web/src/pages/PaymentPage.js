import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./PaymentPage.css";

const PaymentPage = () => {
    const navigate = useNavigate();
    const [products, setProducts] = useState([]);
    const [deliveryMethod, setDeliveryMethod] = useState("pickup");
    const [paymentMethod, setPaymentMethod] = useState("pay_on_delivery");
    const [userAddress, setUserAddress] = useState(null);
    const [cardDetails, setCardDetails] = useState("");

    useEffect(() => {
        const selectedItems = JSON.parse(localStorage.getItem("selectedItems")) || [];

        if (selectedItems.length > 0) {
            fetchProducts(selectedItems);
        }

        const userId = localStorage.getItem("userID");
        if (userId) {
            fetchUserAddress(userId);
        }
    }, []);

    const fetchProducts = async (ids) => {
        try {
            const token = localStorage.getItem("token");
            const fetchedProducts = [];

            for (const id of ids) {
                const response = await fetch(`http://localhost:8080/products/${id}`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${token}`,
                    },
                });

                if (!response.ok) {
                    console.error(`Failed to fetch product with ID ${id}`);
                    continue;
                }

                const productData = await response.json();
                fetchedProducts.push(productData);
            }

            setProducts(fetchedProducts);
        } catch (error) {
            console.error("Error fetching products:", error);
        }
    };

    const fetchUserAddress = async (userId) => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch(`http://localhost:8080/users/${userId}/address`, {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (!response.ok) throw new Error("Failed to fetch user address");

            const data = await response.json();
            setUserAddress(data.address);
        } catch (error) {
            console.error("Error fetching user address:", error);
            setUserAddress("Address not found");
        }
    };


    const handleBack = () => {
        navigate(-1);
    };

    const handlePaymentMethodChange = (method) => {
        setPaymentMethod(method);
    };

    const handleCardDetailsChange = (event) => {
        setCardDetails(event.target.value);
    };

    const handleOrder = async () => {
        const orderDetails = {
            products: products.map((product) => ({
                product_id: product.ID,
                quantity: 1,
                price: product.price
            })),
            deliveryMethod,
            paymentMethod,
            userAddress,
            totalPrice: products.reduce((acc, product) => acc + product.price, 0),
        };

        console.log("Sending order data:", orderDetails);

        try {
            const userId = parseInt(localStorage.getItem("userID"));
            if (isNaN(userId)) {
                throw new Error("Invalid user ID");
            }

            const response = await fetch("http://localhost:8080/orders", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${localStorage.getItem("token")}`,
                },
                body: JSON.stringify({
                    user_id: userId,
                    delivery_method: deliveryMethod,
                    address: userAddress,
                    total_price: orderDetails.totalPrice,
                    order_items: orderDetails.products,
                }),
            });

            const responseText = await response.text();
            console.log("Server response:", responseText);

            if (!response.ok) {
                throw new Error(`Failed to create order: ${responseText}`);
            }

            const responseData = JSON.parse(responseText);
            console.log("Order placed:", responseData);

            await Promise.all(
                orderDetails.products.map(async (item) => {
                    // Получение ID корзины для каждого продукта
                    const cartResponse = await fetch(
                        `http://localhost:8080/cart/${userId}/${item.product_id}`
                    );

                    const cartResponseText = await cartResponse.text();
                    console.log(`Cart response for product ${item.product_id}:`, cartResponseText);

                    if (!cartResponse.ok) {
                        console.error(`Failed to get cart ID for product ${item.product_id}`);
                        return;
                    }

                    let cartData;
                    try {
                        cartData = JSON.parse(cartResponseText);
                        console.log(`Parsed cart data for product ${item.product_id}:`, cartData);
                    } catch (error) {
                        console.error("Error parsing cart response:", error);
                        return;
                    }

                    const cartId = cartData.cart_id; // Используем cart_id вместо id
                    console.log(`Cart ID for product ${item.product_id}:`, cartId);

                    // Проверка типа данных cartId
                    console.log("Type of cartId:", typeof cartId);

                    if (!cartId) {
                        console.error(`No cartId found for product ${item.product_id}`);
                        return;
                    }

                    const deleteResponse = await fetch(
                        `http://localhost:8080/cart/${cartId}`, // Правильный путь с cartId
                        {
                            method: "DELETE",
                            headers: {
                                Authorization: `Bearer ${localStorage.getItem("token")}`,
                            },
                        }
                    );

                    if (!deleteResponse.ok) {
                        console.error("Failed to delete item from cart");
                    } else {
                        console.log(`Deleted item with cartId ${cartId}`);
                    }
                })
            );


            alert("Order placed successfully!");
            navigate("/cart");
        } catch (error) {
            console.error("Error placing order:", error);
            alert("There was an error placing your order.");
        }
    };





    return (
        <div className="payment-page">
            <h2>Payment</h2>
            {userAddress ? (
                <div className="user-address">
                    <h3>Shipping Address: {userAddress}</h3>
                </div>
            ) : (
                <p>Loading address...</p>
            )}

            {products.length > 0 ? (
                <div className="selected-items">
                    <h3>You have selected the following items:</h3>
                    <div className="product-list">
                        {products.map((product, index) => (
                            <div key={index} className="product-card">
                                <h3>{product.name}</h3>
                                <p>Price: ${product.price}</p>
                            </div>
                        ))}
                    </div>
                </div>
            ) : (
                <p>No items selected for payment.</p>
            )}

            <div className="delivery-method">
                <h3>Select a delivery method</h3>
                <div className="delivery-option">
                    <input
                        type="radio"
                        id="pickup"
                        name="delivery"
                        value="pickup"
                        checked={deliveryMethod === "pickup"}
                        onChange={() => setDeliveryMethod("pickup")}
                    />
                    <label htmlFor="pickup">Pickup</label>
                </div>
                <div className="delivery-option">
                    <input
                        type="radio"
                        id="courier"
                        name="delivery"
                        value="courier"
                        checked={deliveryMethod === "courier"}
                        onChange={() => setDeliveryMethod("courier")}
                    />
                    <label htmlFor="courier">Courier</label>
                </div>
            </div>

            <div className="payment-method">
                <h3>Select a payment method</h3>
                <div
                    className={`payment-option ${paymentMethod === "pay_now" ? "selected" : ""}`}
                    onClick={() => handlePaymentMethodChange("pay_now")}
                >
                    Pay Now
                </div>
                <div
                    className={`payment-option ${paymentMethod === "pay_on_delivery" ? "selected" : ""}`}
                    onClick={() => handlePaymentMethodChange("pay_on_delivery")}
                >
                    Pay on Delivery
                </div>
            </div>

            {paymentMethod === "pay_now" && (
                <div className="card-details">
                    <div className="input-group">
                        <input
                            type="text"
                            name="number"
                            value={cardDetails.number}
                            onChange={handleCardDetailsChange}
                            placeholder="Card Number"
                        />
                        <input
                            type="text"
                            name="expiry"
                            value={cardDetails.expiry}
                            onChange={handleCardDetailsChange}
                            className="short-input"
                            placeholder="MM/YY"
                        />
                        <input
                            type="text"
                            name="cvc"
                            value={cardDetails.cvc}
                            onChange={handleCardDetailsChange}
                            className="short-input"
                            placeholder="CVC"
                        />
                    </div>
                    <input
                        type="text"
                        name="cardholder"
                        value={cardDetails.cardholder}
                        onChange={handleCardDetailsChange}
                        placeholder="Cardholder Name"
                    />
                </div>
            )}

            <button className="order-button" onClick={handleOrder}>
                Place Order
            </button>

            <button className="back-button" onClick={handleBack}>
                Back
            </button>
        </div>
    );
};

export default PaymentPage;
