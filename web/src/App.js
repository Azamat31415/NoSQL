import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import ProductsPage from "./pages/ProductsPage";
import Home from "./pages/Home";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Profile from "./pages/Profile";
import "./App.css";
import CartPage from "./pages/CartPage";
import PaymentPage from "./pages/PaymentPage";
import OrderHistory from "./pages/OrderHistory";
import AdminPanel from "./pages/AdminPanel";
import EditProduct from "./pages/EditProduct";
import AddPet from "./pages/AddPet";
import SubscriptionPage from "./pages/SubscriptionPage";
import SubscriptionPaymentPage from "./pages/SubscriptionPaymentPage";
import EditPet from "./pages/EditPet";


function App() {
    return (
        <Router>
            <Navbar />
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/products/:category" element={<ProductsPage />} />
                <Route path="/products/:category/:subcategory" element={<ProductsPage />} />
                <Route path="/products/:category/:subcategory/:type" element={<ProductsPage />} />
                <Route path="/register" element={<Register />} />
                <Route path="/login" element={<Login />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/cart" element={<CartPage />} />
                <Route path="/payment" element={<PaymentPage />} />
                <Route path="/order-history" element={<OrderHistory />} />
                <Route path="/admin-panel" element={<AdminPanel />} />
                <Route path="/edit-product/:id" element={<EditProduct />} />
                <Route path="/add-pet" element={<AddPet />} />
                <Route path="/subscription" element={<SubscriptionPage />} />
                <Route path="/subpayment" element={<SubscriptionPaymentPage />} />
                <Route path="/edit-pet/:id" element={<EditPet />} />
            </Routes>
        </Router>
    );
}

export default App;
