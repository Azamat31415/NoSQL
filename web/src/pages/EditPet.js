import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import "./Pets.css";

const EditPet = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: "",
        species: "",
        age: "",
    });

    const token = localStorage.getItem("token");

    useEffect(() => {
        if (!token) {
            navigate("/login");
            return;
        }

        fetch(`http://localhost:8080/pets/${id}`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${token}`,
            },
        })
            .then((response) => response.json())
            .then((data) => setFormData({
                name: data.Name,
                species: data.Species,
                age: data.Age,
            }))
            .catch((error) => console.error("Error fetching pet data:", error));
    }, [id, navigate, token]);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const updatedPet = {
            name: formData.name,
            species: formData.species,
            age: parseInt(formData.age) || 0,
        };

        try {
            const response = await fetch(`http://localhost:8080/pets/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(updatedPet),
            });

            if (!response.ok) {
                throw new Error("Failed to update pet");
            }

            alert("Pet updated successfully!");
            navigate("/profile");
        } catch (error) {
            console.error("Error:", error);
            alert("Error updating pet");
        }
    };

    return (
        <div className="edit-pet-container">
            <h2>Edit Pet</h2>
            <form className="edit-pet-form" onSubmit={handleSubmit}>
                <label>
                    Name:
                    <input
                        type="text"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        required
                    />
                </label>
                <label>
                    Species:
                    <select name="species" value={formData.species} onChange={handleChange} required>
                        <option value="Dog">Dog</option>
                        <option value="Cat">Cat</option>
                        <option value="Bird">Bird</option>
                        <option value="Rodent">Rodent</option>
                        <option value="Reptile">Reptile</option>
                    </select>
                </label>
                <label>
                    Age:
                    <input
                        type="number"
                        name="age"
                        value={formData.age}
                        onChange={handleChange}
                        min="0"
                    />
                </label>
                <button type="submit" className="save-pet-button">Save Changes</button>
            </form>
        </div>
    );
};

export default EditPet;
