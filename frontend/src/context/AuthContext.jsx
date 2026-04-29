import { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const role = localStorage.getItem('role');
        if (token) {
            setUser({ token, role });
        }
        setLoading(false);
    }, []);

    const login = (tokenData, role) => {
        // Handle both string token and object with access_token
        const token = typeof tokenData === 'string' ? tokenData : tokenData?.access_token;
        localStorage.setItem('token', token);
        if (role) localStorage.setItem('role', role);
        setUser({ token, role });
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        localStorage.removeItem('cart');
        setUser(null);
    };

    const isAuthenticated = !!user?.token;
    const isAdmin = user?.role === 'admin';

    return (
        <AuthContext.Provider value={{ user, login, logout, isAuthenticated, isAdmin, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) throw new Error('useAuth must be used within an AuthProvider');
    return context;
};

export default AuthContext;
