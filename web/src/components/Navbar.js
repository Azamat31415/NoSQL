import { useState, useCallback, useEffect } from "react";
import { Link } from "react-router-dom";

const categories = [
    {
        name: "Dog",
        subcategories: [
            { title: "Feed", items: ["Dry", "Wet", "Super premium class", "Grain-free", "Hypoallergenic"] },
            { title: "Treats", items: ["Chews", "Snacks", "Biscuits", "Training treats", "Dental sticks"] },
            { title: "Care and hygiene", items: ["Shampoos", "Claw clippers", "Ear cleaners", "Flea and tick protection"] },
            { title: "Toilets and trays", items: ["Puppy pads", "Waste bags", "Odor neutralizers"] },
            { title: "Clothes and shoes", items: ["Winter coats", "Raincoats", "Boots", "Bandanas"] },
            { title: "Ammunition", items: ["Collars", "Leashes", "Muzzles", "Harnesses", "ID tags"] },
            { title: "Beds and houses", items: ["Orthopedic beds", "Crates", "Travel carriers"] },
            { title: "Toys", items: ["Chew toys", "Rope toys", "Interactive toys", "Balls", "Squeaky toys"] },
        ],
    },
    {
        name: "Cat",
        subcategories: [
            { title: "Feed", items: ["Dry", "Wet", "Super premium class", "Grain-free", "Hairball control"] },
            { title: "Treats", items: ["Crunchy treats", "Soft treats", "Catnip-infused treats"] },
            { title: "Care and hygiene", items: ["Shampoos", "Nail clippers", "Flea and tick treatment", "Ear cleaners"] },
            { title: "Toilets and trays", items: ["Litter boxes", "Clumping litter", "Silica gel litter"] },
            { title: "Scratching posts", items: ["Simple scratching posts", "Cat trees", "Wall-mounted scratchers"] },
            { title: "Beds", items: ["Heated beds", "Covered beds", "Window perches"] },
            { title: "Toys", items: ["Balls", "Mice", "Feather wands", "Interactive laser toys"] },
        ],
    },
    {
        name: "Rodent",
        subcategories: [
            { title: "Feed", items: ["Pellets", "Seed mix", "Hay", "Vitamin supplements"] },
            { title: "Treats", items: ["Dried fruits", "Nuts", "Crunchy sticks"] },
            { title: "Cages and accessories", items: ["Cages", "Tunnels", "Exercise wheels", "Hammocks"] },
            { title: "Bedding", items: ["Wood shavings", "Paper bedding", "Straw"] },
            { title: "Care", items: ["Brushes", "Nail trimmers", "Tooth care"] },
        ],
    },
    {
        name: "Bird",
        subcategories: [
            { title: "Feed", items: ["Grain mix", "Pellets", "Fruit blend", "Egg food"] },
            { title: "Treats", items: ["Seed sticks", "Fruit treats", "Mineral blocks"] },
            { title: "Cages and accessories", items: ["Bird cages", "Perches", "Nest boxes"] },
            { title: "Toys", items: ["Mirrors", "Swings", "Chewable toys"] },
            { title: "Care", items: ["Beak and claw care", "Feather sprays"] },
        ],
    },
    {
        name: "Reptile",
        subcategories: [
            { title: "Feed", items: ["Live insects", "Frozen mice", "Pellets", "Dried food"] },
            { title: "Terrariums and equipment", items: ["Terrariums", "Heat lamps", "UVB bulbs", "Humidity control"] },
            { title: "Substrate and bedding", items: ["Coconut fiber", "Sand", "Moss"] },
            { title: "Decor", items: ["Hiding spots", "Climbing branches", "Water dishes"] },
            { title: "Care", items: ["Calcium supplements", "Shedding aids", "Vitamin sprays"] },
        ],
    },
    {
        name: "Vet. Pharmacy",
        subcategories: [
            { title: "Supplements", items: ["Multivitamins", "Joint health", "Digestive support"] },
            { title: "Medicine", items: ["Antibiotics", "Pain relievers", "Deworming tablets"] },
            { title: "First aid", items: ["Wound disinfectants", "Bandages", "Eye and ear drops"] },
            { title: "Parasite control", items: ["Flea and tick treatments", "Deworming solutions"] },
        ],
    },
];

const Navbar = () => {
    const [activeCategory, setActiveCategory] = useState(null);
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (token) {
            setIsLoggedIn(true);
            fetch("http://localhost:8080/profile", {
                method: "GET",
                headers: { Authorization: `Bearer ${token}` },
            })
                .then((response) => response.json())
                .catch((error) => console.error("Error fetching profile:", error));
        }
    }, []);


    useEffect(() => {
        const handleClickOutside = (event) => {
            if (!event.target.closest(".dropdown")) {
                setActiveCategory(null);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);


    const handleCategoryClick = useCallback((categoryName) => {
        setActiveCategory((prev) => (prev === categoryName ? null : categoryName));
    }, []);

    const handleOutsideClick = (event) => {
        if (!event.target.closest(".dropdown")) {
            setActiveCategory(null);
        }
    };

    return (
        <>
            <header className="header">
                <Link to="/" className="logo">Pet Store</Link>
                <div className="nav-links">
                    <Link to="/order-history">Order History</Link>
                    <Link to="/subscription">Subscription</Link>
                    {!isLoggedIn && <><Link to="/register">Register</Link><Link to="/login">Login</Link></>}
                    {isLoggedIn && <><Link to="/profile">Profile</Link></>}
                    <Link to="/cart">Cart</Link>
                </div>
            </header>

            <nav onClick={handleOutsideClick}>
                <ul className="nav-list">
                    {categories.map((category) => (
                        <li key={category.name} className="dropdown">
                            <button
                                className="dropbtn"
                                onClick={() => handleCategoryClick(category.name)}
                            >
                                {category.name}
                            </button>
                            {activeCategory === category.name && (
                                <div className="dropdown-content">
                                    <table className="category-table">
                                        <thead>
                                        <tr>
                                            {category.subcategories.map((subcategory) => (
                                                <th key={subcategory.title}>{subcategory.title}</th>
                                            ))}
                                        </tr>
                                        </thead>
                                        <tbody>
                                        <tr>
                                            {category.subcategories.map((subcategory) => (
                                                <td key={subcategory.title}>
                                                    {subcategory.items.map((item) => (
                                                        <div key={item} className="subcategory-item">
                                                            <Link
                                                                to={`/products/${category.name.toLowerCase()}/${subcategory.title.toLowerCase()}/${item.toLowerCase()}`}
                                                            >
                                                                {item}
                                                            </Link>
                                                        </div>
                                                    ))}
                                                </td>
                                            ))}
                                        </tr>
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </li>
                    ))}
                </ul>
            </nav>
        </>
    );
};

export default Navbar;
