import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import ProductCard from "../components/ProductCard";

const ProductPage = () => {
    const { category, subcategory, type } = useParams();
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                let url = `http://localhost:8080/products?category=${category}`;
                if (subcategory) url += `&subcategory=${subcategory}`;
                if (type) url += `&type=${type}`;

                const response = await fetch(url);
                if (!response.ok) throw new Error("Failed to fetch products");

                const data = await response.json();
                setProducts(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [category, subcategory, type]);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error: {error}</p>;

    return (
        <div>
            <h2>Products for {category} - {subcategory} - {type}</h2>
            <div className="product-list">
                {products.length > 0 ? (
                    products.map((product) => (
                        <ProductCard key={product.id} name={product.name} description={product.description} price={product.price} image={product.image} />
                    ))
                ) : (
                    <p>No products found</p>
                )}
            </div>
        </div>
    );
};

export default ProductPage;
