import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Pets.css";

const AddPet = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: "",
        species: "",
        age: "",
    });

    const token = localStorage.getItem("token");
    const userID = localStorage.getItem("userID");

    if (!token || !userID) {
        navigate("/login");
        return null;
    }

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const petData = {
            ...formData,
            age: parseInt(formData.age) || 0,
            userID: parseInt(userID),
        };

        console.log("Sending pet data:", petData);

        try {
            const response = await fetch("http://localhost:8080/pets", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(petData),
            });

            if (!response.ok) {
                throw new Error("Failed to add pet");
            }

            alert("Pet added successfully!");
            navigate("/profile");
        } catch (error) {
            console.error("Error:", error);
            alert("Error adding pet");
        }
    };

    return (
        <div className="add-pet-container">
            <h2>Add a New Pet</h2>
            <form className="add-pet-form" onSubmit={handleSubmit}>
                <label>
                    Name:
                    <input type="text" name="name" value={formData.name} onChange={handleChange} required />
                </label>
                <label>
                    Species:
                    <select name="species" value={formData.species} onChange={handleChange}>
                        <option value="Dog">Dog</option>
                        <option value="Cat">Cat</option>
                        <option value="Bird">Bird</option>
                        <option value="Rodent">Rodent</option>
                        <option value="Reptile">Reptile</option>
                    </select>
                </label>
                <label>
                    Age:
                    <input type="number" name="age" value={formData.age} onChange={handleChange} min="0"/>
                </label>
                <button type="submit" className="add-pet-button">
                    Add Pet
                </button>
            </form>
        </div>
    );
};

export default AddPet;
