import React, { useState, useEffect } from "react";
import "./OrderHistory.css";

const OrderHistory = () => {
    const [orders, setOrders] = useState([]);

    useEffect(() => {
        const userId = localStorage.getItem("userID");
        console.log("Fetched userID:", userId);

        if (!userId) {
            console.error("User ID is missing");
            return;
        }

        fetch(`http://localhost:8080/order-history/${userId}`)
            .then((res) => {
                if (!res.ok) {
                    throw new Error(`HTTP error! Status: ${res.status}`);
                }
                return res.json();
            })
            .then((data) => setOrders(data))
            .catch((err) => console.error("Error fetching orders:", err));
    }, []);

    return (
        <div className="order-history-container">
            {orders.length > 0 ? (
                orders.map((order) => (
                    <div key={order.ID} className="order-card">
                        <h3>Order #{order.ID}</h3>
                        <p className="order-details">
                            <strong>Status:</strong>{" "}
                            <span className={`status-${order.Status.toLowerCase()}`}>
                                {order.Status}
                            </span>
                        </p>
                        <p className="order-details"><strong>Delivery:</strong> {order.DeliveryMethod}</p>
                        <p className="order-details"><strong>Address:</strong> {order.Address || "N/A"}</p>
                        <p className="order-details"><strong>Total Price:</strong> ${order.TotalPrice.toFixed(2)}</p>
                        <div>
                            <strong>Items:</strong>
                            {order.OrderItems.length > 0 ? (
                                order.OrderItems.map((item) => (
                                    <p key={item.ID} className="order-item">
                                        Product ID: {item.ProductID} x {item.Quantity} (${item.Price.toFixed(2)})
                                    </p>
                                ))
                            ) : (
                                <p>No items in this order.</p>
                            )}
                        </div>
                    </div>
                ))
            ) : (
                <p>No orders found.</p>
            )}
        </div>
    );
};

export default OrderHistory;
