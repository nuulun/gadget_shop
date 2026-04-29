import { useEffect, useState } from 'react';
import { Plus, Edit, Trash2, Search, Loader2 } from 'lucide-react';
import { getProducts, deleteProduct } from '../api';
import { Link } from 'react-router-dom';

const Admin = () => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        loadProducts();
    }, []);

    const loadProducts = async () => {
        try {
            const { data } = await getProducts();
            setProducts(data);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id) => {
        if (window.confirm('Are you sure you want to delete this product?')) {
            try {
                await deleteProduct(id);
                setProducts(products.filter(p => (p._id || p.id) !== id));
            } catch (error) {
                alert('Failed to delete product');
            }
        }
    };

    const filtered = products.filter(p =>
        p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        p.brand.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) return <div className="flex justify-center p-10"><Loader2 className="animate-spin text-primary" /></div>;

    return (
        <div className="min-h-screen bg-background p-8">
            <div className="max-w-7xl mx-auto">
                <div className="flex justify-between items-center mb-8">
                    <div>
                        <h1 className="text-3xl font-bold text-slate-800">Admin Dashboard</h1>
                        <p className="text-secondary">Manage your inventory</p>
                    </div>
                    <Link to="/admin/new" className="bg-primary hover:bg-blue-600 text-white px-6 py-2 rounded-xl font-bold shadow-glow flex items-center gap-2 transition-all">
                        <Plus size={20} /> Add Product
                    </Link>
                </div>

                <div className="bg-surface rounded-3xl shadow-soft border border-slate-100 overflow-hidden">
                    {/* Toolbar */}
                    <div className="p-4 border-b border-slate-100 flex gap-4">
                        <div className="relative flex-1 max-w-sm">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                            <input
                                type="text"
                                placeholder="Search products..."
                                className="w-full pl-10 pr-4 py-2 rounded-lg border border-slate-200 focus:outline-none focus:border-primary transition-colors"
                                value={searchTerm}
                                onChange={e => setSearchTerm(e.target.value)}
                            />
                        </div>
                    </div>

                    {/* Table */}
                    <div className="overflow-x-auto">
                        <table className="w-full text-left border-collapse">
                            <thead>
                                <tr className="bg-slate-50 text-slate-600 text-sm uppercase tracking-wider">
                                    <th className="p-4 font-semibold">Product</th>
                                    <th className="p-4 font-semibold">Category</th>
                                    <th className="p-4 font-semibold">Price</th>
                                    <th className="p-4 font-semibold">Stock</th>
                                    <th className="p-4 font-semibold text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-50">
                                {filtered.map(product => (
                                    <tr key={product._id || product.id} className="hover:bg-slate-50 transition-colors">
                                        <td className="p-4">
                                            <div className="flex items-center gap-3">
                                                <div className="w-10 h-10 bg-slate-100 rounded-lg flex items-center justify-center text-lg">
                                                    💻
                                                </div>
                                                <div>
                                                    <div className="font-bold text-slate-800">{product.name}</div>
                                                    <div className="text-xs text-slate-500">{product.brand}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="p-4">
                                            <span className="bg-blue-50 text-primary px-2 py-1 rounded text-xs font-bold uppercase">
                                                {product.category}
                                            </span>
                                        </td>
                                        <td className="p-4 font-bold text-slate-700">${product.price}</td>
                                        <td className="p-4">
                                            <span className={`px-2 py-1 rounded text-xs font-bold ${product.stock < 5 ? 'bg-red-50 text-red-500' : 'bg-green-50 text-green-600'}`}>
                                                {product.stock} in stock
                                            </span>
                                        </td>
                                        <td className="p-4 text-right">
                                            <div className="flex items-center justify-end gap-2">
                                                <Link to={`/admin/edit/${product._id || product.id}`} className="p-2 text-slate-400 hover:text-primary hover:bg-blue-50 rounded-lg transition-colors flex items-center justify-center">
                                                    <Edit size={18} />
                                                </Link>
                                                <button
                                                    onClick={() => handleDelete(product._id || product.id)}
                                                    className="p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                                                >
                                                    <Trash2 size={18} />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        {filtered.length === 0 && (
                            <div className="p-8 text-center text-secondary">No products found.</div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Admin;
