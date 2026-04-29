import { useState } from 'react';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import { Trash2, Plus, Minus, CreditCard, ShoppingBag, Banknote } from 'lucide-react';
import { createOrder } from '../api';
import { Link, useNavigate } from 'react-router-dom';

const Cart = () => {
    const { cart, removeFromCart, updateQuantity, clearCart, total } = useCart();
    const { isAuthenticated } = useAuth();
    const navigate = useNavigate();
    const [paymentMethod, setPaymentMethod] = useState('card');
    const [loading, setLoading] = useState(false);

    const handleCheckout = async () => {
        if (!isAuthenticated) {
            alert('Please login to complete your order.');
            navigate('/login');
            return;
        }
        setLoading(true);
        try {
            // Map cart items to the format our order-service expects
            const orderData = {
                items: cart.map(item => ({
                    product_id: item.id,          // our service uses snake_case
                    quantity: item.quantity,
                })),
            };
            await createOrder(orderData);
            alert('Order placed successfully!');
            clearCart();
            navigate('/');
        } catch (error) {
            const msg = error.response?.data?.error || error.message || 'Backend unreachable';
            alert(`Failed to place order: ${msg}`);
            if (error.response?.status === 401) navigate('/login');
        } finally {
            setLoading(false);
        }
    };

    if (cart.length === 0) {
        return (
            <div className="min-h-screen bg-background flex flex-col items-center justify-center p-8 text-center">
                <div className="bg-surface p-6 rounded-full shadow-soft mb-6">
                    <ShoppingBag size={48} className="text-secondary" />
                </div>
                <h2 className="text-2xl font-bold text-slate-800 mb-2">Your cart is empty</h2>
                <p className="text-secondary mb-6">Browse the market to find something you like.</p>
                <Link to="/market" className="bg-primary text-white px-6 py-2 rounded-xl font-medium shadow-glow hover:bg-blue-600 transition-colors">
                    Browse Market
                </Link>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background p-8">
            <div className="max-w-4xl mx-auto">
                <h1 className="text-3xl font-extrabold text-slate-800 mb-8">Your Cart</h1>

                <div className="space-y-4 mb-8">
                    {cart.map(item => (
                        <div key={item.id} className="bg-surface rounded-2xl p-4 flex items-center gap-4 shadow-soft border border-slate-100">
                            <img src={item.image} alt={item.name}
                                className="w-20 h-20 object-cover rounded-xl bg-slate-100"
                                onError={e => { e.target.src = 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=80'; }} />
                            <div className="flex-1">
                                <h3 className="font-bold text-slate-800">{item.name}</h3>
                                <p className="text-sm text-secondary capitalize">{item.category}</p>
                                <p className="text-primary font-bold">${item.price}</p>
                            </div>
                            <div className="flex items-center gap-2">
                                <button onClick={() => updateQuantity(item.id, item.quantity - 1)}
                                    className="p-1 rounded-lg bg-slate-100 hover:bg-slate-200 transition-colors">
                                    <Minus size={16} />
                                </button>
                                <span className="w-8 text-center font-bold">{item.quantity}</span>
                                <button onClick={() => updateQuantity(item.id, item.quantity + 1)}
                                    className="p-1 rounded-lg bg-slate-100 hover:bg-slate-200 transition-colors">
                                    <Plus size={16} />
                                </button>
                            </div>
                            <p className="font-bold text-slate-800 w-20 text-right">
                                ${(item.price * item.quantity).toFixed(2)}
                            </p>
                            <button onClick={() => removeFromCart(item.id)}
                                className="p-2 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                                <Trash2 size={18} />
                            </button>
                        </div>
                    ))}
                </div>

                {/* Payment method */}
                <div className="bg-surface rounded-2xl p-6 shadow-soft border border-slate-100 mb-6">
                    <h3 className="font-bold text-slate-800 mb-4">Payment Method</h3>
                    <div className="flex gap-3">
                        {[{value:'card',label:'Card',icon:<CreditCard size={18}/>},{value:'cash',label:'Cash',icon:<Banknote size={18}/>}].map(m => (
                            <button key={m.value} onClick={() => setPaymentMethod(m.value)}
                                className={`flex items-center gap-2 px-4 py-2 rounded-xl border-2 font-medium transition-all ${paymentMethod === m.value ? 'border-primary bg-blue-50 text-primary' : 'border-slate-200 text-secondary hover:border-primary/40'}`}>
                                {m.icon}{m.label}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Summary */}
                <div className="bg-surface rounded-2xl p-6 shadow-soft border border-slate-100">
                    <div className="flex justify-between items-center mb-4">
                        <span className="text-secondary">Subtotal</span>
                        <span className="font-bold">${total.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between items-center mb-6 text-xl font-extrabold">
                        <span>Total</span>
                        <span className="text-primary">${total.toFixed(2)}</span>
                    </div>
                    <button onClick={handleCheckout} disabled={loading}
                        className="w-full bg-primary hover:bg-blue-600 text-white font-bold py-3 rounded-xl shadow-glow flex items-center justify-center gap-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                        {loading ? 'Placing Order…' : <><CreditCard size={20}/> Place Order</>}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Cart;
