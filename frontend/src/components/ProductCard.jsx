import { ShoppingCart } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useCart } from '../context/CartContext';

const FALLBACK = 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=400';

const ProductCard = ({ product }) => {
    const { addToCart } = useCart();

    const handleAddToCart = (e) => {
        e.preventDefault();
        e.stopPropagation();
        addToCart(product);
    };

    // Support both field names (image from new service, image_url from old)
    const imgSrc = product.image || product.image_url || FALLBACK;

    return (
        <div className="bg-surface rounded-2xl p-4 shadow-soft transition-all duration-300 hover:-translate-y-1 hover:shadow-xl border border-slate-100 group relative flex flex-col">
            {/* Image */}
            <div className="h-44 bg-slate-50 rounded-xl mb-4 relative overflow-hidden flex items-center justify-center">
                <Link to={`/product/${product.id}`} className="absolute inset-0">
                    <img
                        src={imgSrc}
                        alt={product.name}
                        className="w-full h-full object-cover transition-transform group-hover:scale-105"
                        onError={e => { e.target.src = FALLBACK; }}
                    />
                </Link>

                {/* Add to Cart button */}
                <button
                    onClick={handleAddToCart}
                    className="absolute bottom-3 right-3 bg-primary text-white p-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity shadow-glow z-10 hover:scale-110 active:scale-95"
                    title="Add to Cart"
                >
                    <ShoppingCart size={18} />
                </button>

                {/* Out of stock badge */}
                {product.stock === 0 && (
                    <span className="absolute top-2 left-2 bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded-full z-10">
                        Out of stock
                    </span>
                )}
            </div>

            {/* Content */}
            <div className="flex flex-col flex-1 space-y-2">
                <div className="flex justify-between items-center">
                    <span className="text-xs font-bold text-primary bg-blue-50 px-2 py-0.5 rounded-full uppercase tracking-wider">
                        {product.category}
                    </span>
                    {product.brand && (
                        <span className="text-xs font-medium text-slate-400">{product.brand}</span>
                    )}
                </div>

                <Link to={`/product/${product.id}`} className="flex-1">
                    <h3 className="font-bold text-slate-800 leading-tight line-clamp-2 hover:text-primary transition-colors">
                        {product.name}
                    </h3>
                </Link>

                <div className="flex justify-between items-center pt-1">
                    <span className="text-xl font-bold text-slate-900">
                        ${Number(product.price).toFixed(2)}
                    </span>
                    {product.stock > 0 && product.stock < 5 && (
                        <span className="text-xs text-amber-500 font-semibold">
                            Only {product.stock} left!
                        </span>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ProductCard;
