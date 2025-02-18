import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../App.css";

const Register = () => {
    const [formData, setFormData] = useState({
        email: "",
        password: "",
        first_name: "",
        last_name: "",
        address: "",
        phone: "",
    });

    const [error, setError] = useState("");
    const [errorMessage, setErrorMessage] = useState("");
    const navigate = useNavigate();

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setErrorMessage("");

        const response = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(formData),
        });

        try {
            if (!response.ok) {
                const textResponse = await response.text();
                setErrorMessage("User with this email already exists");
                console.error("Server response:", textResponse);
                return;
            }

            const data = await response.json();

            if (data.message && data.message.includes("duplicate key") && data.message.includes("already exists")) {
                setErrorMessage("User with this email already exists.");
            } else {
                navigate("/login");
            }
        } catch (error) {
            setErrorMessage("Error parsing response. Please try again.");
            console.error("Error parsing JSON:", error);
        }
    };

    return (
        <div className="auth-container">
            <h2>Register</h2>
            {error && <p className="error">{error}</p>}
            {errorMessage && <p className="error">{errorMessage}</p>}
            <form onSubmit={handleSubmit}>
                <input type="email" name="email" placeholder="Email" onChange={handleChange} required />
                <input type="password" name="password" placeholder="Password" onChange={handleChange} required />
                <input type="text" name="first_name" placeholder="First Name" onChange={handleChange} required />
                <input type="text" name="last_name" placeholder="Last Name" onChange={handleChange} />
                <input type="text" name="address" placeholder="Address" onChange={handleChange} />
                <input type="text" name="phone" placeholder="Phone" onChange={handleChange} />
                <button type="submit">Register</button>
            </form>
        </div>
    );
};

export default Register;
