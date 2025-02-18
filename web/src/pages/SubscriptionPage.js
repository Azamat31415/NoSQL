import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "./subscryp.css";

const foodOptions = [
    { type: "Dry", price: 100 },
    { type: "Wet", price: 120 },
    { type: "Grain-free", price: 110 },
];

const SubscriptionPage = () => {
    const [selectedFood, setSelectedFood] = useState(foodOptions[0].type);
    const [price, setPrice] = useState(foodOptions[0].price);
    const [subscription, setSubscription] = useState(null);
    const navigate = useNavigate();

    const fetchSubscription = async () => {
        const token = localStorage.getItem("token");
        const userId = localStorage.getItem("userID");

        if (!token || !userId) {
            console.log("No token or userID found.");
            return;
        }

        try {
            const response = await fetch(`http://localhost:8080/subscriptions/${userId}`, {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (response.ok) {
                const data = await response.json();
                console.log("Subscription data received:", data);

                let formattedDate = "Unknown";
                if (data.RenewalDate) {
                    const parsedDate = new Date(data.RenewalDate);
                    if (!isNaN(parsedDate.getTime())) {
                        formattedDate = parsedDate.toLocaleDateString();
                    }
                }

                setSubscription({
                    id: data.ID,
                    type: data.Type,
                    renewalDate: formattedDate,
                });
            } else {
                console.warn("No active subscription found.");
                setSubscription(null);
            }
        } catch (error) {
            console.error("Error fetching subscription:", error);
        }
    };

    useEffect(() => {
        fetchSubscription();
    }, []);

    useEffect(() => {
        const selected = foodOptions.find(food => food.type === selectedFood);
        if (selected) {
            setPrice(selected.price);
        }
    }, [selectedFood]);

    const handleSubscription = async () => {
        const token = localStorage.getItem("token");
        const userId = localStorage.getItem("userID");

        if (!token || !userId) {
            alert("You need to log in to subscribe.");
            navigate("/login");
            return;
        }

        const subscriptionData = {
            user_id: parseInt(userId),
            interval_days: 90,
            type: selectedFood,
            status: "active",
        };

        try {
            const response = await fetch("http://localhost:8080/subscriptions", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(subscriptionData),
            });

            if (response.ok) {
                alert("Subscription successful!");
                await fetchSubscription();
                navigate("/subpayment", { state: { price, selectedFood } });
            } else {
                const responseText = await response.text();
                alert(`Failed to subscribe: ${responseText}`);
            }
        } catch (error) {
            alert("An error occurred. Please check your connection.");
        }
    };

    const handleCancelSubscription = async () => {
        if (!subscription || !subscription.id) {
            console.log("No subscription to cancel.");
            return;
        }

        const token = localStorage.getItem("token");

        try {
            console.log("Sending DELETE request to:", `http://localhost:8080/subscriptions/${subscription.id}`);
            const response = await fetch(`http://localhost:8080/subscriptions/${subscription.id}`, {
                method: "DELETE",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (response.ok) {
                alert("Subscription canceled.");
                console.log("Subscription successfully deleted.");
                setSubscription(null);
                await fetchSubscription(); // Обновляем данные после удаления
            } else {
                const responseText = await response.text();
                alert(`Failed to cancel subscription: ${responseText}`);
                console.log("Response error:", responseText);
            }
        } catch (error) {
            console.error("Error canceling subscription:", error);
            alert("An error occurred while canceling the subscription.");
        }
    };

    return (
        <div className="subscription-container">
            <h2>3-Month Pet Food Subscription</h2>
            {subscription ? (
                <>
                    <p>Your subscription for <strong>{subscription.type}</strong> is active until {subscription.renewalDate}.</p>
                    <button className="cancel-button" onClick={handleCancelSubscription}>
                        Cancel Subscription
                    </button>
                </>
            ) : (
                <>
                    <p>Receive your selected pet food every month for 3 months.</p>
                    <label>
                        Choose food type:
                        <select
                            value={selectedFood}
                            onChange={(e) => setSelectedFood(e.target.value)}
                        >
                            {foodOptions.map((food) => (
                                <option key={food.type} value={food.type}>{food.type}</option>
                            ))}
                        </select>
                    </label>
                    <p className="price">Price: {price} $</p>
                    <button className="subscribe-button" onClick={handleSubscription}>
                        Subscribe Now
                    </button>
                </>
            )}
        </div>
    );
};

export default SubscriptionPage;
