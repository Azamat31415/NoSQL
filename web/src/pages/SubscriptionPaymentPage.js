import React, { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import "./subpayment.css";

const SubscriptionPaymentPage = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [paymentMethod, setPaymentMethod] = useState("pay_now");
    const [cardDetails, setCardDetails] = useState({ number: "", expiry: "", cvc: "", cardholder: "" });
    const [price, setPrice] = useState(0);
    const [selectedFood, setSelectedFood] = useState("");
    const [errors, setErrors] = useState({});

    useEffect(() => {
        if (location.state?.price && location.state?.selectedFood) {
            setPrice(location.state.price);
            setSelectedFood(location.state.selectedFood);
        }
    }, [location.state]);

    const handlePaymentMethodChange = (method) => {
        setPaymentMethod(method);
    };

    const handleCardDetailsChange = (event) => {
        setCardDetails({ ...cardDetails, [event.target.name]: event.target.value });
    };

    const validateFields = () => {
        let newErrors = {};
        if (paymentMethod === "pay_now") {
            if (!cardDetails.number.trim()) newErrors.number = "Card number is required";
            if (!cardDetails.expiry.trim()) newErrors.expiry = "Expiry date is required";
            if (!cardDetails.cvc.trim()) newErrors.cvc = "CVC is required";
            if (!cardDetails.cardholder.trim()) newErrors.cardholder = "Cardholder name is required";
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handlePayment = async () => {
        if (paymentMethod === "pay_now" && !validateFields()) {
            return;
        }

        const token = localStorage.getItem("token");
        const userId = localStorage.getItem("userID");

        if (!token) {
            alert("You need to log in to proceed.");
            navigate("/login");
            return;
        }

        if (!userId) {
            alert("User ID is missing. Please log in again.");
            navigate("/login");
            return;
        }

        const paymentData = {
            user_id: parseInt(userId),
            amount: price,
            payment_method: paymentMethod,
            status: "completed",
        };

        try {
            const response = await fetch("http://localhost:8080/subpayment", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(paymentData),
            });

            if (!response.ok) {
                const responseText = await response.text();
                alert(`Payment failed: ${responseText}`);
                return;
            }

            // Устанавливаем дату окончания подписки (через 3 месяца от текущей даты)
            const endDate = new Date();
            endDate.setMonth(endDate.getMonth() + 3);

            // Сохраняем подписку в localStorage
            localStorage.setItem("subscription", JSON.stringify({
                type: selectedFood,
                endDate: endDate.toISOString(),
            }));

            alert("Subscription payment successful!");
            navigate("/subscription");
        } catch (error) {
            alert("There was an error processing your payment.");
        }
    };

    return (
        <div className="subpayment">
            <h2>Subscription Payment</h2>
            <p>Price: ${price}</p>
            <div className="payment-method">
                <h3>Select a payment method</h3>
                <div
                    className={`payment-option ${paymentMethod === "pay_now" ? "selected" : ""}`}
                    onClick={() => handlePaymentMethodChange("pay_now")}
                >
                    Pay Now
                </div>
            </div>

            {paymentMethod === "pay_now" && (
                <div className="card-details">
                    <input
                        type="text"
                        name="number"
                        value={cardDetails.number}
                        onChange={handleCardDetailsChange}
                        placeholder="Card Number"
                    />
                    {errors.number && <p className="error">{errors.number}</p>}

                    <input
                        type="text"
                        name="expiry"
                        value={cardDetails.expiry}
                        onChange={handleCardDetailsChange}
                        className="short-input"
                        placeholder="MM/YY"
                    />
                    {errors.expiry && <p className="error">{errors.expiry}</p>}

                    <input
                        type="text"
                        name="cvc"
                        value={cardDetails.cvc}
                        onChange={handleCardDetailsChange}
                        className="short-input"
                        placeholder="CVC"
                    />
                    {errors.cvc && <p className="error">{errors.cvc}</p>}

                    <input
                        type="text"
                        name="cardholder"
                        value={cardDetails.cardholder}
                        onChange={handleCardDetailsChange}
                        placeholder="Cardholder Name"
                    />
                    {errors.cardholder && <p className="error">{errors.cardholder}</p>}
                </div>
            )}

            <button className="pay-button" onClick={handlePayment}>
                Confirm Payment
            </button>
            <button className="back-button" onClick={() => navigate(-1)}>
                Back
            </button>
        </div>
    );
};

export default SubscriptionPaymentPage;
