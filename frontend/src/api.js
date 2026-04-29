import axios from 'axios';

const API = axios.create({
    baseURL: '/api'
});

// Attach token to every request
API.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => Promise.reject(error));

// Auth
export const login = (credentials) => API.post('/auth/login', {
    login: credentials.email || credentials.login,
    password: credentials.password
});

export const register = (userData) => API.post('/auth/register', {
    login: userData.login || userData.email.split('@')[0],
    email: userData.email,
    password: userData.password,
    first_name: userData.fullName ? userData.fullName.split(' ')[0] : (userData.first_name || ''),
    last_name: userData.fullName ? (userData.fullName.split(' ')[1] || '') : (userData.last_name || ''),
    phone: userData.phone || '',
    age: userData.age || 0,
});

export const getMe = () => API.get('/auth/me');

// Products
export const getProducts = () => API.get('/products');
export const getProduct = (id) => API.get(`/products/${id}`);
export const createProduct = (data) => API.post('/products', data);
export const updateProduct = (id, data) => API.put(`/products/${id}`, data);
export const deleteProduct = (id) => API.delete(`/products/${id}`);

// Orders
export const createOrder = (orderData) => API.post('/orders', orderData);
export const getMyOrders = () => API.get('/orders/my-orders');

// Cart (client-side localStorage)
export const getCart = () => {
    const cart = localStorage.getItem('cart');
    return Promise.resolve({ data: cart ? JSON.parse(cart) : [] });
};
export const updateCart = (items) => {
    localStorage.setItem('cart', JSON.stringify(items));
    return Promise.resolve({ data: items });
};

const api = {
    login, register, getMe,
    getProducts, getProduct, createProduct, updateProduct, deleteProduct,
    createOrder, getMyOrders,
    getCart, updateCart,
};

export default api;
