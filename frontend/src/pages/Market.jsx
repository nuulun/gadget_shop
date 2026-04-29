import { useEffect, useState } from 'react';
import { Loader2, AlertCircle, Search } from 'lucide-react';
import ProductCard from '../components/ProductCard';
import { getProducts } from '../api';

const CATEGORIES = ['all', 'keyboard', 'mouse', 'monitor', 'laptop', 'headset', 'webcam'];

const Market = () => {
    const [allProducts, setAllProducts] = useState([]);
    const [filtered, setFiltered] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [category, setCategory] = useState('all');
    const [search, setSearch] = useState('');

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                const response = await getProducts();
                // Handle both {items:[...]} and plain [...] responses
                const data = response.data;
                const list = Array.isArray(data) ? data : (data.items || data.products || []);
                console.log('[Market] products loaded:', list.length);
                setAllProducts(list);
                setFiltered(list);
            } catch (err) {
                console.error('[Market] failed to fetch products:', err);
                setError('Failed to load products. Is the backend running?');
            } finally {
                setLoading(false);
            }
        };
        fetchProducts();
    }, []);

    // Filter whenever category or search changes
    useEffect(() => {
        let result = allProducts;
        if (category !== 'all') {
            result = result.filter(p => p.category === category);
        }
        if (search.trim()) {
            const q = search.toLowerCase();
            result = result.filter(p =>
                p.name?.toLowerCase().includes(q) ||
                p.brand?.toLowerCase().includes(q) ||
                p.description?.toLowerCase().includes(q)
            );
        }
        setFiltered(result);
    }, [category, search, allProducts]);

    if (loading) {
        return (
            <div className="min-h-screen bg-background flex items-center justify-center">
                <Loader2 className="animate-spin text-primary" size={48} />
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen bg-background flex flex-col items-center justify-center text-center p-8">
                <div className="bg-red-50 p-4 rounded-full mb-4">
                    <AlertCircle className="text-red-500" size={32} />
                </div>
                <h2 className="text-2xl font-bold text-slate-800 mb-2">Connection Error</h2>
                <p className="text-secondary mb-6">{error}</p>
                <button
                    onClick={() => window.location.reload()}
                    className="bg-primary text-white px-6 py-2 rounded-xl font-medium shadow-glow hover:bg-blue-600 transition-colors"
                >
                    Retry
                </button>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background">
            {/* Header */}
            <div className="max-w-7xl mx-auto px-8 pt-10 pb-6">
                <h1 className="text-4xl font-extrabold text-slate-800 tracking-tight mb-1">
                    Gadget<span className="text-primary">Market</span>
                </h1>
                <p className="text-secondary text-lg mb-6">Premium tech for developers.</p>

                {/* Search + Filters */}
                <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
                    {/* Search */}
                    <div className="relative flex-1 max-w-md">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                        <input
                            type="text"
                            placeholder="Search products…"
                            value={search}
                            onChange={e => setSearch(e.target.value)}
                            className="w-full pl-10 pr-4 py-2.5 bg-surface border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all"
                        />
                    </div>

                    {/* Category pills */}
                    <div className="flex gap-2 flex-wrap">
                        {CATEGORIES.map(cat => (
                            <button
                                key={cat}
                                onClick={() => setCategory(cat)}
                                className={`px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wide transition-all ${
                                    category === cat
                                        ? 'bg-primary text-white shadow-glow'
                                        : 'bg-surface text-secondary border border-slate-200 hover:border-primary/50 hover:text-primary'
                                }`}
                            >
                                {cat}
                            </button>
                        ))}
                    </div>
                </div>
            </div>

            {/* Product Grid */}
            <div className="max-w-7xl mx-auto px-8 pb-16">
                {filtered.length === 0 ? (
                    <div className="text-center py-24 text-secondary">
                        <p className="text-xl font-medium mb-2">No products found</p>
                        <p className="text-sm">
                            {allProducts.length === 0
                                ? 'The product database is empty. Check that product-service started correctly.'
                                : 'Try a different search or category.'}
                        </p>
                    </div>
                ) : (
                    <>
                        <p className="text-sm text-secondary mb-6">
                            Showing {filtered.length} of {allProducts.length} products
                        </p>
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                            {filtered.map(p => (
                                <ProductCard key={p.id} product={p} />
                            ))}
                        </div>
                    </>
                )}
            </div>
        </div>
    );
};

export default Market;
