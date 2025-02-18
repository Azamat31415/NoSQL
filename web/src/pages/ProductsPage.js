import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import ProductCard from "../components/ProductCard";

const ProductsPage = () => {
    const { category, subcategory, type } = useParams();
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch(`http://localhost:8080/products?category=${category}&subcategory=${subcategory}&type=${type}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP Error ${response.status}`);
                }
                return response.json();
            })
            .then(data => setProducts(data))
            .catch(error => setError(error.message))
            .finally(() => setLoading(false));
    }, [category, subcategory, type]);

    const addToCart = (productId) => {
        const token = localStorage.getItem("token");
        const userId = localStorage.getItem("user_id");

        if (!token || !userId) {
            alert("You need to log in to add items to the cart");
            return;
        }

        const cartItem = { user_id: userId, product_id: productId, quantity: 1 };

        console.log("Adding to cart:", cartItem);

        fetch("http://localhost:8080/cart", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(cartItem),
        })
            .then((response) => {
                console.log("Response status:", response.status);
                return response.json();
            })
            .then((json) => {
                console.log("Server response:", json);
                alert("Product added to cart successfully");
            })
            .catch((error) => alert(error.message));
    };



    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error: {error}</p>;

    return (
        <div className="products-page">
            <h2>Products: {category} - {subcategory} - {type}</h2>
            <div className="product-list">
                {products.length > 0 ? (
                    products.map((product) => (
                        <ProductCard key={product.ID} product={product} addToCart={addToCart} />
                    ))
                ) : (
                    <p>No products found</p>
                )}
            </div>
        </div>
    );
};

export default ProductsPage;
