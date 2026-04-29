import { ArrowRight, ShoppingBag } from 'lucide-react';
import { Link } from 'react-router-dom';

const Home = () => {
    return (
        <div className="flex flex-col items-center justify-center min-h-[calc(100vh-64px)] text-center px-4">
            <div className="max-w-3xl space-y-8 animate-fade-in">
                <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 text-primary font-medium text-sm">
                    <ShoppingBag size={16} />
                    <span>The Future of Tech is Here</span>
                </div>

                <h1 className="text-5xl md:text-7xl font-black tracking-tight text-slate-900 leading-tight">
                    Discover the <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-blue-600">Extraordinary</span>
                </h1>

                <p className="text-xl text-slate-600 max-w-2xl mx-auto leading-relaxed">
                    Experience the latest in gadget innovation. Premium quality, curated selection, and seamless shopping experience defined for the modern age.
                </p>

                <div className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4">
                    <Link
                        to="/market"
                        className="group relative px-8 py-4 bg-primary text-white font-bold rounded-2xl shadow-glow hover:bg-blue-600 transition-all flex items-center gap-2 overflow-hidden"
                    >
                        <span className="relative z-10 flex items-center gap-2">
                            Start Shopping <ArrowRight className="group-hover:translate-x-1 transition-transform" />
                        </span>
                    </Link>

                    <Link
                        to="/register"
                        className="px-8 py-4 bg-white text-slate-700 font-bold rounded-2xl border border-slate-200 hover:bg-slate-50 transition-all shadow-sm"
                    >
                        Create Account
                    </Link>
                </div>
            </div>

            {/* Background Decor Elements */}
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-primary/5 rounded-full blur-3xl -z-10" />
            <div className="absolute bottom-0 right-0 w-[300px] h-[300px] bg-blue-400/10 rounded-full blur-3xl -z-10" />
        </div>
    );
};

export default Home;
