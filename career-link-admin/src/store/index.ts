// src/store/index.ts
import { configureStore } from '@reduxjs/toolkit';
import authReducer from './slices/authSlice';
// Import other reducers as needed

export const store = configureStore({
    reducer: {
        auth: authReducer,
        // Add other reducers here
    },
});

// Export types for TypeScript usage
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;