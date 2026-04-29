import { useEffect, useState } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { ArrowLeft, Save, Plus, Trash2, Loader2 } from 'lucide-react';
import { getProduct, createProduct, updateProduct } from '../api';

const AdminProductForm = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const isEditing = !!id;
    const [loading, setLoading] = useState(isEditing);

    const [formData, setFormData] = useState({
        name: '',
        brand: '',
        price: '',
        category: '',
        stock: '',
        specs: {}
    });

    const [specList, setSpecList] = useState([]); // [{key: '', value: ''}]

    useEffect(() => {
        if (isEditing) {
            loadProduct();
        }
    }, [id]);

    const loadProduct = async () => {
        try {
            const { data } = await getProduct(id);
            setFormData({
                name: data.name,
                brand: data.brand,
                price: data.price,
                category: data.category,
                stock: data.stock,
                specs: data.specs || {}
            });
            setSpecList(Object.entries(data.specs || {}).map(([key, value]) => ({ key, value })));
        } catch (error) {
            console.error(error);
            alert("Failed to load product");
        } finally {
            setLoading(false);
        }
    };

    const handleSpecChange = (index, field, value) => {
        const newSpecs = [...specList];
        newSpecs[index][field] = value;
        setSpecList(newSpecs);
    };

    const addSpec = () => setSpecList([...specList, { key: '', value: '' }]);
    const removeSpec = (index) => setSpecList(specList.filter((_, i) => i !== index));

    const handleSubmit = async (e) => {
        e.preventDefault();
        const specsObj = specList.reduce((acc, curr) => {
            if (curr.key) acc[curr.key] = curr.value;
            return acc;
        }, {});

        const payload = {
            ...formData,
            price: Number(formData.price),
            stock: Number(formData.stock),
            specs: specsObj
        };

        try {
            if (isEditing) {
                await updateProduct(id, payload);
            } else {
                await createProduct(payload);
            }
            navigate('/admin');
        } catch (error) {
            console.error(error);
            alert("Failed to save product");
        }
    };

    if (loading) return <div className="flex justify-center p-10"><Loader2 className="animate-spin" /></div>;

    return (
        <div className="min-h-screen bg-background p-8">
            <div className="max-w-3xl mx-auto">
                <Link to="/admin" className="flex items-center text-secondary hover:text-primary mb-6 transition-colors font-medium">
                    <ArrowLeft size={20} className="mr-2" /> Back to Dashboard
                </Link>

                <div className="bg-surface rounded-3xl shadow-soft p-8 border border-slate-100">
                    <h1 className="text-2xl font-bold text-slate-800 mb-6">{isEditing ? 'Edit Product' : 'New Product'}</h1>

                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">Product Name</label>
                                <input
                                    type="text" required
                                    className="w-full bg-background border border-slate-200 rounded-xl py-3 px-4 focus:outline-none focus:border-primary transition-colors"
                                    value={formData.name}
                                    onChange={e => setFormData({ ...formData, name: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">Brand</label>
                                <input
                                    type="text" required
                                    className="w-full bg-background border border-slate-200 rounded-xl py-3 px-4 focus:outline-none focus:border-primary transition-colors"
                                    value={formData.brand}
                                    onChange={e => setFormData({ ...formData, brand: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">Category</label>
                                <input
                                    type="text" required
                                    className="w-full bg-background border border-slate-200 rounded-xl py-3 px-4 focus:outline-none focus:border-primary transition-colors"
                                    value={formData.category}
                                    onChange={e => setFormData({ ...formData, category: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">Price ($)</label>
                                <input
                                    type="number" required min="0" step="0.01"
                                    className="w-full bg-background border border-slate-200 rounded-xl py-3 px-4 focus:outline-none focus:border-primary transition-colors"
                                    value={formData.price}
                                    onChange={e => setFormData({ ...formData, price: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-slate-700 mb-2">Stock Quantity</label>
                                <input
                                    type="number" required min="0"
                                    className="w-full bg-background border border-slate-200 rounded-xl py-3 px-4 focus:outline-none focus:border-primary transition-colors"
                                    value={formData.stock}
                                    onChange={e => setFormData({ ...formData, stock: e.target.value })}
                                />
                            </div>
                        </div>

                        {/* Specs Editor */}
                        <div className="border-t border-slate-100 pt-6">
                            <div className="flex justify-between items-center mb-4">
                                <h3 className="font-bold text-slate-800">Specifications</h3>
                                <button type="button" onClick={addSpec} className="text-primary text-sm font-bold flex items-center gap-1 hover:underline">
                                    <Plus size={16} /> Add Spec
                                </button>
                            </div>
                            <div className="space-y-3">
                                {specList.map((spec, index) => (
                                    <div key={index} className="flex gap-3">
                                        <input
                                            type="text" placeholder="Key (e.g. RAM)"
                                            className="flex-1 bg-background border border-slate-200 rounded-lg py-2 px-3 focus:outline-none focus:border-primary text-sm transition-colors"
                                            value={spec.key}
                                            onChange={e => handleSpecChange(index, 'key', e.target.value)}
                                        />
                                        <input
                                            type="text" placeholder="Value (e.g. 16GB)"
                                            className="flex-1 bg-background border border-slate-200 rounded-lg py-2 px-3 focus:outline-none focus:border-primary text-sm transition-colors"
                                            value={spec.value}
                                            onChange={e => handleSpecChange(index, 'value', e.target.value)}
                                        />
                                        <button type="button" onClick={() => removeSpec(index)} className="text-slate-400 hover:text-red-500 transition-colors">
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="pt-6 border-t border-slate-100 flex justify-end">
                            <button className="bg-primary hover:bg-blue-600 text-white px-8 py-3 rounded-xl font-bold shadow-glow flex items-center gap-2 transition-all">
                                <Save size={20} /> Save Product
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default AdminProductForm;
