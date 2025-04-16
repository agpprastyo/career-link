// src/api/auth.ts
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const authApi = {
    login: async (credentials: { email: string; password: string }) => {
        const response = await axios.post(`${API_URL}/auth/login`, credentials);
        return response.data;
    },

    getCurrentUser: async () => {
        const token = localStorage.getItem('token');
        if (!token) return null;

        const response = await axios.get(`${API_URL}/users/me`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        return response.data;
    }
};