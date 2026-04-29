import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Loader2, AlertCircle, ShoppingCart, ArrowLeft, Plus, Minus } from 'lucide-react';
import { getProduct } from '../api';
import { useCart } from '../context/CartContext';

const FALLBACK = 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=600';

const ProductDetail = () => {
    const { id } = useParams();
    const { addToCart } = useCart();
    const [product, setProduct] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [quantity, setQuantity] = useState(1);
    const [added, setAdded] = useState(false);

    useEffect(() => {
        if (!id) return;
        setLoading(true);
        getProduct(id)
            .then(res => setProduct(res.data))
            .catch(err => {
                console.error('[ProductDetail] fetch error:', err);
                setError('Could not load product details.');
            })
            .finally(() => setLoading(false));
    }, [id]);

    const handleAddToCart = () => {
        for (let i = 0; i < quantity; i++) addToCart(product);
        setAdded(true);
        setTimeout(() => setAdded(false), 2000);
    };

    if (loading) return (
        <div className="min-h-screen bg-background flex items-center justify-center">
            <Loader2 className="animate-spin text-primary" size={48} />
        </div>
    );

    if (error || !product) return (
        <div className="min-h-screen bg-background flex flex-col items-center justify-center text-center p-8">
            <AlertCircle className="text-red-400 mb-3" size={40} />
            <p className="text-slate-700 font-medium mb-4">{error || 'Product not found.'}</p>
            <Link to="/market" className="text-primary font-bold hover:underline flex items-center gap-1">
                <ArrowLeft size={16} /> Back to Market
            </Link>
        </div>
    );

    const imgSrc = product.image || product.image_url || FALLBACK;
    const maxQty = product.stock || 99;

    return (
        <div className="min-h-screen bg-background p-8">
            <div className="max-w-6xl mx-auto">
                <Link to="/market" className="inline-flex items-center gap-2 text-secondary hover:text-primary mb-8 transition-colors">
                    <ArrowLeft size={18} /> Back to Market
                </Link>

                <div className="bg-surface rounded-3xl shadow-soft p-8 grid grid-cols-1 md:grid-cols-2 gap-12 border border-slate-100">
                    {/* Image */}
                    <div className="bg-slate-50 rounded-2xl overflow-hidden min-h-[380px] flex items-center justify-center">
                        <img src={imgSrc} alt={product.name}
                            className="w-full h-full object-cover"
                            onError={e => { e.target.src = FALLBACK; }} />
                    </div>

                    {/* Info */}
                    <div className="flex flex-col">
                        <div className="mb-3">
                            <span className="text-xs font-bold text-primary bg-blue-50 px-3 py-1 rounded-full uppercase tracking-wider">
                                {product.category}
                            </span>
                        </div>
                        <h1 className="text-3xl font-extrabold text-slate-800 mb-1 leading-tight">{product.name}</h1>
                        {product.brand && <p className="text-slate-500 font-medium mb-4">{product.brand}</p>}
                        {product.description && (
                            <p className="text-slate-600 text-sm mb-6 leading-relaxed">{product.description}</p>
                        )}

                        <div className="text-4xl font-bold text-slate-900 mb-6">
                            ${Number(product.price).toFixed(2)}
                        </div>

                        {/* Stock */}
                        <p className={`text-sm font-medium mb-4 ${product.stock === 0 ? 'text-red-500' : product.stock < 5 ? 'text-amber-500' : 'text-green-600'}`}>
                            {product.stock === 0 ? 'Out of stock' : product.stock < 5 ? `Only ${product.stock} left!` : `${product.stock} in stock`}
                        </p>

                        {/* Quantity */}
                        {product.stock > 0 && (
                            <div className="flex items-center gap-4 mb-6">
                                <span className="text-slate-600 font-medium">Quantity:</span>
                                <div className="flex items-center gap-3 bg-slate-100 rounded-xl p-1.5">
                                    <button onClick={() => setQuantity(q => Math.max(1, q - 1))}
                                        className="p-1.5 rounded-lg hover:bg-slate-200 transition-colors disabled:opacity-40"
                                        disabled={quantity <= 1}>
                                        <Minus size={16} />
                                    </button>
                                    <span className="font-bold text-lg w-8 text-center">{quantity}</span>
                                    <button onClick={() => setQuantity(q => Math.min(maxQty, q + 1))}
                                        className="p-1.5 rounded-lg hover:bg-slate-200 transition-colors disabled:opacity-40"
                                        disabled={quantity >= maxQty}>
                                        <Plus size={16} />
                                    </button>
                                </div>
                            </div>
                        )}

                        <button
                            onClick={handleAddToCart}
                            disabled={product.stock === 0}
                            className="w-full bg-primary hover:bg-blue-600 text-white font-bold py-4 rounded-2xl shadow-glow flex items-center justify-center gap-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed active:scale-95"
                        >
                            <ShoppingCart size={22} />
                            {added ? '✓ Added to Cart!' : `Add${quantity > 1 ? ` ${quantity}` : ''} to Cart`}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProductDetail;
