import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import "./edit.css";

const EditProduct = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [product, setProduct] = useState({ name: "", description: "", price: "", image: "" });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch(`http://localhost:8080/products/${id}`)
            .then(response => response.json())
            .then(data => setProduct(data))
            .catch(error => setError("Failed to load product"))
            .finally(() => setLoading(false));
    }, [id]);

    const handleChange = (e) => {
        setProduct({ ...product, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        const token = localStorage.getItem("token");

        fetch(`http://localhost:8080/products/${id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(product),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error("Failed to update product");
                }
                return response.json();
            })
            .then(() => {
                alert("Product updated successfully");
                navigate("/admin-panel");
            })
            .catch(error => alert(error.message));
    };

    if (loading) return <p>Loading...</p>;
    if (error) return <p>{error}</p>;

    return (
        <div className="edit-product-container">
            <h2>Edit Product</h2>
            <form onSubmit={handleSubmit} className="edit-product-form">
                <label>
                    Name:
                    <input type="text" name="name" value={product.name} onChange={handleChange} required />
                </label>
                <label>
                    Description:
                    <textarea name="description" value={product.description} onChange={handleChange} required />
                </label>
                <label>
                    Price:
                    <input type="number" name="price" value={product.price} onChange={handleChange} required />
                </label>
                <button type="submit">Save Changes</button>
                <button type="button" onClick={() => navigate("/admin-panel")}>Cancel</button>
            </form>
        </div>
    );
};

export default EditProduct;