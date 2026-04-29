import { useState } from 'react';
import { Mail, Lock, ArrowRight, Loader2, AlertCircle } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { login as apiLogin } from '../api';
import { useAuth } from '../context/AuthContext';

const Login = () => {
    const navigate = useNavigate();
    const { login } = useAuth();
    const [formData, setFormData] = useState({ email: '', password: '' });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
        setError(null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        try {
            const response = await apiLogin(formData);
            const data = response.data;
            // data is { access_token, refresh_token }
            login(data.access_token);
            navigate('/');
        } catch (err) {
            setError(err.response?.data?.error || 'Invalid credentials');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-[calc(100vh-64px)]">
            <div className="bg-surface p-10 rounded-3xl shadow-soft border border-slate-100 w-full max-w-md mx-4">
                <div className="mb-8 text-center">
                    <h2 className="text-3xl font-extrabold text-slate-800">Welcome back</h2>
                    <p className="text-secondary mt-2">Enter your credentials</p>
                </div>
                {error && (
                    <div className="mb-4 p-3 bg-red-50 text-red-500 rounded-xl flex items-center gap-2 text-sm">
                        <AlertCircle size={16} /> {error}
                    </div>
                )}
                <form className="space-y-6" onSubmit={handleSubmit}>
                    <div>
                        <label className="block text-sm font-semibold text-slate-700 mb-2">Email or Login</label>
                        <div className="relative">
                            <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={20} />
                            <input type="text" name="email" required value={formData.email}
                                onChange={handleChange} placeholder="nurlan@example.kz"
                                className="w-full bg-background border border-slate-200 rounded-xl py-3 pl-11 pr-4 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all" />
                        </div>
                    </div>
                    <div>
                        <label className="block text-sm font-semibold text-slate-700 mb-2">Password</label>
                        <div className="relative">
                            <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={20} />
                            <input type="password" name="password" required value={formData.password}
                                onChange={handleChange} placeholder="••••••••"
                                className="w-full bg-background border border-slate-200 rounded-xl py-3 pl-11 pr-4 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all" />
                        </div>
                    </div>
                    <button disabled={loading}
                        className="w-full bg-primary hover:bg-blue-600 text-white font-bold py-3 rounded-xl shadow-glow flex items-center justify-center gap-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                        {loading ? <Loader2 className="animate-spin" /> : <>Sign In <ArrowRight size={20} /></>}
                    </button>
                </form>
                <p className="text-center mt-8 text-sm text-secondary">
                    Don't have an account? <Link to="/register" className="text-primary font-bold cursor-pointer hover:underline">Register</Link>
                </p>
            </div>
        </div>
    );
};

export default Login;
