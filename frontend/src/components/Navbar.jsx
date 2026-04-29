import { Link, useNavigate } from 'react-router-dom';
import { ShoppingCart, LogIn, LogOut, LayoutDashboard, Store } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

const Navbar = () => {
    const { isAuthenticated, isAdmin, logout } = useAuth();
    const { count } = useCart();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    return (
        <nav className="sticky top-0 z-50 bg-surface/80 backdrop-blur-md border-b border-slate-100 shadow-soft">
            <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
                {/* Logo */}
                <Link to="/" className="text-2xl font-extrabold text-slate-800 tracking-tight">
                    Gadget<span className="text-primary">Shop</span>
                </Link>

                {/* Nav links */}
                <div className="flex items-center gap-2">
                    <Link to="/market"
                        className="flex items-center gap-1.5 px-3 py-2 rounded-xl text-sm font-medium text-secondary hover:text-primary hover:bg-blue-50 transition-all">
                        <Store size={16} /> Market
                    </Link>

                    {isAdmin && (
                        <Link to="/admin"
                            className="flex items-center gap-1.5 px-3 py-2 rounded-xl text-sm font-medium text-secondary hover:text-primary hover:bg-blue-50 transition-all">
                            <LayoutDashboard size={16} /> Admin
                        </Link>
                    )}

                    <Link to="/cart"
                        className="relative flex items-center gap-1.5 px-3 py-2 rounded-xl text-sm font-medium text-secondary hover:text-primary hover:bg-blue-50 transition-all">
                        <ShoppingCart size={16} />
                        Cart
                        {count > 0 && (
                            <span className="absolute -top-1 -right-1 bg-primary text-white text-xs font-bold w-5 h-5 rounded-full flex items-center justify-center">
                                {count > 9 ? '9+' : count}
                            </span>
                        )}
                    </Link>

                    {isAuthenticated ? (
                        <button onClick={handleLogout}
                            className="flex items-center gap-1.5 px-3 py-2 rounded-xl text-sm font-medium text-secondary hover:text-red-500 hover:bg-red-50 transition-all">
                            <LogOut size={16} /> Logout
                        </button>
                    ) : (
                        <Link to="/login"
                            className="flex items-center gap-1.5 px-4 py-2 rounded-xl text-sm font-bold bg-primary text-white shadow-glow hover:bg-blue-600 transition-all">
                            <LogIn size={16} /> Sign In
                        </Link>
                    )}
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
