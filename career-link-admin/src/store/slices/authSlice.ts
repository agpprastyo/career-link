// src/store/slices/authSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { authApi } from '../../api/auth';

export const login = createAsyncThunk(
    'auth/login',
    async (credentials: { email: string; password: string }, { rejectWithValue }) => {
        try {
            const data = await authApi.login(credentials);
            localStorage.setItem('token', data.token);
            return data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data || 'Login failed');
        }
    }
);

export const fetchCurrentUser = createAsyncThunk(
    'auth/fetchCurrentUser',
    async (_, { rejectWithValue }) => {
        try {
            return await authApi.getCurrentUser();
        } catch (error: any) {
            return rejectWithValue(error.response?.data || 'Failed to fetch user');
        }
    }
);

const authSlice = createSlice({
    name: 'auth',
    initialState: {
        user: null,
        token: localStorage.getItem('token'),
        isLoading: false,
        error: null as string | null
    },
    reducers: {
        logout: (state) => {
            state.user = null;
            state.token = null;
            localStorage.removeItem('token');
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(login.pending, (state) => {
                state.isLoading = true;
                state.error = null;
            })
            .addCase(login.fulfilled, (state, action) => {
                state.isLoading = false;
                state.token = action.payload.token;
            })
            .addCase(login.rejected, (state, action) => {
                state.isLoading = false;
                state.error = action.payload as string;
            })
            .addCase(fetchCurrentUser.fulfilled, (state, action) => {
                state.user = action.payload;
            });
    }
});

export const { logout } = authSlice.actions;
export default authSlice.reducer;