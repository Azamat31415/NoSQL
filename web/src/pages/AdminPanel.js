import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./AdminPanel.css";
import ProductCard from "../components/ProductCard";

const AdminPanel = () => {
    const navigate = useNavigate();
    const token = localStorage.getItem("token");
    const role = localStorage.getItem("role");
    const [activeTab, setActiveTab] = useState(null);
    const [orders, setOrders] = useState([]);
    const [productId, setProductId] = useState("");
    const [users, setUsers] = useState([]);
    const [product, setProduct] = useState(null);
    const [subscriptionId, setSubscriptionId] = useState("");
    const [subscription, setSubscription] = useState(null);
    const [productForm, setProductForm] = useState({
        name: "",
        description: "",
        price: "",
        stock: "",
        category: "",
        subcategory: "",
        type: ""
    });

    const categories = {
        "dog": {
            "feed": ["dry", "wet", "super premium class", "grain-free", "hypoallergenic"],
            "treats": ["chews", "snacks", "biscuits", "training treats", "dental sticks"],
            "care and hygiene": ["shampoos", "claw clippers", "ear cleaners", "flea and tick protection"],
            "toilets and trays": ["puppy pads", "waste bags", "odor neutralizers"],
            "clothes and shoes": ["winter coats", "raincoats", "boots", "bandanas"],
            "ammunition": ["collars", "leashes", "muzzles", "harnesses", "id tags"],
            "beds and houses": ["orthopedic beds", "crates", "travel carriers"],
            "toys": ["chew toys", "rope toys", "interactive toys", "balls", "squeaky toys"]
        },
        "cat": {
            "feed": ["dry", "wet", "super premium class", "grain-free", "hairball control"],
            "treats": ["crunchy treats", "soft treats", "catnip-infused treats"],
            "care and hygiene": ["shampoos", "nail clippers", "flea and tick treatment", "ear cleaners"],
            "toilets and trays": ["litter boxes", "clumping litter", "silica gel litter"],
            "scratching posts": ["simple scratching posts", "cat trees", "wall-mounted scratchers"],
            "beds": ["heated beds", "covered beds", "window perches"],
            "toys": ["balls", "mice", "feather wands", "interactive laser toys"]
        },
        "rodent": {
            "feed": ["pellets", "seed mix", "hay", "vitamin supplements"],
            "treats": ["dried fruits", "nuts", "crunchy sticks"],
            "cages and accessories": ["cages", "tunnels", "exercise wheels", "hammocks"],
            "bedding": ["wood shavings", "paper bedding", "straw"],
            "care": ["brushes", "nail trimmers", "tooth care"]
        },
        "bird": {
            "feed": ["grain mix", "pellets", "fruit blend", "egg food"],
            "treats": ["seed sticks", "fruit treats", "mineral blocks"],
            "cages and accessories": ["bird cages", "perches", "nest boxes"],
            "toys": ["mirrors", "swings", "chewable toys"],
            "care": ["beak and claw care", "feather sprays"]
        },
        "reptile": {
            "feed": ["live insects", "frozen mice", "pellets", "dried food"],
            "terrariums and equipment": ["terrariums", "heat lamps", "uvb bulbs", "humidity control"],
            "substrate and bedding": ["coconut fiber", "sand", "moss"],
            "decor": ["hiding spots", "climbing branches", "water dishes"],
            "care": ["calcium supplements", "shedding aids", "vitamin sprays"]
        },
        "vet. pharmacy": {
            "supplements": ["multivitamins", "joint health", "digestive support"],
            "medicine": ["antibiotics", "pain relievers", "deworming tablets"],
            "first aid": ["wound disinfectants", "bandages", "eye and ear drops"],
            "parasite control": ["flea and tick treatments", "deworming solutions"]
        }
    };



    useEffect(() => {
        if (!token || role !== "admin") {
            navigate("/profile");
        }
    }, [navigate, token, role]);

    useEffect(() => {
        if (activeTab === "orders") {
            fetchOrders();
        }
    }, [activeTab]);

    useEffect(() => {
        if (activeTab === "users") {
            fetchUsers();
        }
    }, [activeTab]);


    const fetchUsers = async () => {
        try {
            const response = await fetch("http://localhost:8080/users", {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setUsers(data);
            } else {
                alert("Failed to fetch users");
            }
        } catch (error) {
            console.error("Error fetching users:", error);
        }
    };

    const fetchSubscriptionById = async () => {
        if (!subscriptionId) return;
        try {
            const response = await fetch(`http://localhost:8080/subscriptions/${subscriptionId}`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setSubscription(data);
            } else {
                setSubscription(null);
                alert("Subscription not found");
            }
        } catch (error) {
            console.error("Error fetching subscription:", error);
        }
    };


    const handleSubscriptionIdChange = (e) => {
        setSubscriptionId(e.target.value);
    };

    const fetchProductById = async () => {
        if (!productId) return;
        try {
            const response = await fetch(`http://localhost:8080/products/${productId}`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setProduct(data);
            } else {
                setProduct(null);
                alert("Product not found");
            }
        } catch (error) {
            console.error("Error fetching product:", error);
        }
    };

    const handleProductIdChange = (e) => {
        setProductId(e.target.value);
    };



    const fetchOrders = async () => {
        try {
            const response = await fetch("http://localhost:8080/orders", {
                headers: { Authorization: `Bearer ${token}` }
            });
            const data = await response.json();
            setOrders(data);
        } catch (error) {
            console.error("Error fetching orders:", error);
        }
    };

    const handleProductChange = (e) => {
        const { name, value } = e.target;

        if (name === "category") {
            setProductForm({
                ...productForm,
                category: value,
                subcategory: "",
                type: ""
            });
        } else if (name === "subcategory") {
            setProductForm({
                ...productForm,
                subcategory: value,
                type: ""
            });
        } else {
            setProductForm({ ...productForm, [name]: value });
        }
    };

    const handleProductSubmit = async (e) => {
        e.preventDefault();
        const productData = {
            ...productForm,
            price: parseFloat(productForm.price),
            stock: parseInt(productForm.stock, 10)
        };
        try {
            const response = await fetch("http://localhost:8080/products", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify(productData)
            });
            if (response.ok) {
                alert("Product added successfully");
                setProductForm({
                    name: "",
                    description: "",
                    price: "",
                    stock: "",
                    category: "",
                    subcategory: "",
                    type: ""
                });
            } else {
                alert("Failed to add product");
            }
        } catch (error) {
            console.error("Error adding product:", error);
        }
    };

    const handleStatusChange = async (orderId, newStatus) => {
        try {
            const response = await fetch(`http://localhost:8080/orders/${orderId}/status/update`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({ status: newStatus }),
            });

            if (response.ok) {
                alert("Order status updated successfully");
                setOrders(prevOrders =>
                    prevOrders.map(order =>
                        order.ID === orderId ? { ...order, Status: newStatus } : order
                    )
                );
            }else {
                alert("Failed to update order status");
            }
        } catch (error) {
            console.error("Error updating order status:", error);
        }
    };

    const orderStatuses = ["Pending", "Shipped", "Delivered", "Cancelled", "Returned",];


    const renderContent = () => {
        switch (activeTab) {
            case "orders":
                return (
                    <div className="admin-content spaced-content">
                        <h3>Orders</h3>
                        <div className="orders-grid">
                            {orders.map(order => (
                                <div key={order.ID} className="order-card">
                                    <p><strong>Order ID:</strong> {order.ID}</p>
                                    <p><strong>User ID:</strong> {order.UserID}</p>
                                    <p><strong>Delivery:</strong> {order.DeliveryMethod}</p>
                                    <p><strong>Address:</strong> {order.Address}</p>
                                    <p><strong>Status:</strong> {order.Status}</p>
                                    <p><strong>Total:</strong> ${order.TotalPrice.toFixed(2)}</p>

                                    <select
                                        onChange={(e) => handleStatusChange(order.ID, e.target.value.toLowerCase())}
                                        defaultValue={order.Status.toLowerCase()}
                                    >
                                        {orderStatuses.map(status => (
                                            <option key={status} value={status}>{status}</option>
                                        ))}
                                    </select>
                                </div>
                            ))}
                        </div>
                    </div>
                );
            case "products":
                return (
                    <div className="admin-content spaced-content">
                        <h3>Add Product</h3>
                        <form onSubmit={handleProductSubmit}>
                            <input name="name" placeholder="Name" onChange={handleProductChange} required />
                            <input name="description" placeholder="Description" onChange={handleProductChange} required />
                            <input name="price" type="number" placeholder="Price" onChange={handleProductChange} required />
                            <input name="stock" type="number" placeholder="Stock" onChange={handleProductChange} required />
                            <select name="category" onChange={handleProductChange} required>
                                <option value="">Select Category</option>
                                {Object.keys(categories).map(cat => (
                                    <option key={cat} value={cat}>{cat}</option>
                                ))}
                            </select>
                            {productForm.category && (
                                <select name="subcategory" onChange={handleProductChange} required>
                                    <option value="">Select Subcategory</option>
                                    {Object.keys(categories[productForm.category]).map(sub => (
                                        <option key={sub} value={sub}>{sub}</option>
                                    ))}
                                </select>
                            )}
                            {productForm.subcategory && (
                                <select name="type" onChange={handleProductChange} required>
                                    <option value="">Select Type</option>
                                    {categories[productForm.category][productForm.subcategory].map(type => (
                                        <option key={type} value={type}>{type}</option>
                                    ))}
                                </select>
                            )}
                            <button type="submit">Add Product</button>
                        </form>
                        <h3>Find Product by ID</h3>
                        <input type="text" placeholder="Enter Product ID" value={productId} onChange={handleProductIdChange} />
                        <button onClick={fetchProductById}>Find Product</button>
                        {product && <ProductCard product={product} />}
                    </div>
                );
            case "subscriptions":
                return (
                    <div className="admin-content spaced-content">
                        <h3>Find Subscription by ID</h3>
                        <input
                            type="text"
                            placeholder="Enter Subscription ID"
                            value={subscriptionId}
                            onChange={(e) => setSubscriptionId(e.target.value)}
                        />
                        <button onClick={fetchSubscriptionById}>Find Subscription</button>

                        {subscription && (
                            <div className="subscription-card">
                                <p><strong>ID:</strong> {subscription.ID}</p>
                                <p><strong>User ID:</strong> {subscription.UserID}</p>
                                <p><strong>Type:</strong> {subscription.Type}</p>
                                <p><strong>Status:</strong> {subscription.Status}</p>
                                <p><strong>Start Date:</strong> {new Date(subscription.StartDate).toLocaleDateString()}</p>
                                <p><strong>Renewal Date:</strong> {new Date(subscription.RenewalDate).toLocaleDateString()}</p>
                                <p><strong>Interval (Days):</strong> {subscription.IntervalDays}</p>
                            </div>
                        )}
                    </div>
                );
            case "users":
                return (
                    <div className="admin-content spaced-content">
                        <h3>Users List</h3>
                        <table className="user-table">
                            <thead>
                            <tr>
                                <th>ID</th>
                                <th>Email</th>
                                <th>First Name</th>
                                <th>Last Name</th>
                                <th>Phone</th>
                                <th>Address</th>
                                <th>Role</th>
                                <th>Created At</th>
                            </tr>
                            </thead>
                            <tbody>
                            {users.map(user => (
                                <tr key={user.ID}>
                                    <td>{user.ID}</td>
                                    <td>{user.Email}</td>
                                    <td>{user.FirstName}</td>
                                    <td>{user.LastName}</td>
                                    <td>{user.Phone}</td>
                                    <td>{user.Address}</td>
                                    <td>{user.Role}</td>
                                    <td>{new Date(user.CreatedAt).toLocaleDateString()}</td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
                );
            default:
                return <div className="admin-content spaced-content">Select an option above</div>;
        }
    };

    return (
        <div className="admin-panel">
            <h2 className="admin-title">Admin Panel</h2>
            <div className="admin-buttons spaced sticky-buttons">
                <button onClick={() => setActiveTab("orders")}>Orders</button>
                <button onClick={() => setActiveTab("subscriptions")}>Subscriptions</button>
                <button onClick={() => setActiveTab("products")}>Products</button>
                <button onClick={() => setActiveTab("users")}>Users</button>
            </div>
            <div className="content-wrapper">{renderContent()}</div>
        </div>
    );
};

export default AdminPanel;
