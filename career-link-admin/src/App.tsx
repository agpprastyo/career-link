// src/App.tsx
import {JSX, useEffect} from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { CircularProgress, Box, CssBaseline } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import { useAppDispatch, useAppSelector } from './hooks/redux';
import { fetchCurrentUser } from './store/slices/authSlice';

// Layouts
import AdminLayout from './components/layout/AdminLayout';

// Pages
import LoginPage from './features/auth/LoginPage';
import DashboardPage from './features/dashboard/DashboardPage';
// import UsersListPage from './features/users/UsersListPage';
// import UserDetailPage from './features/users/UserDetailPage';
// import JobsListPage from './features/jobs/JobsListPage';
import NotFoundPage from './components/NotFoundPage';

// Theme
import theme from './theme';

// Protected Route Component
const ProtectedRoute = ({ children }: { children: JSX.Element }) => {
    const { token, user, isLoading } = useAppSelector(state => state.auth);

    if (isLoading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                <CircularProgress />
            </Box>
        );
    }

    if (!token || !user) {
        return <Navigate to="/login" replace />;
    }

    return children;
};

function App() {
    const dispatch = useAppDispatch();
    const { token } = useAppSelector(state => state.auth);

    useEffect(() => {
        if (token) {
            dispatch(fetchCurrentUser());
        }
    }, [dispatch, token]);

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<LoginPage />} />

                    {/* Protected Admin Routes */}
                    <Route
                        path="/"
                        element={
                            <ProtectedRoute>
                                <AdminLayout />
                            </ProtectedRoute>
                        }
                    >
                        <Route index element={<Navigate to="/dashboard" replace />} />
                        <Route path="dashboard" element={<DashboardPage />} />
                        {/*<Route path="users" element={<UsersListPage />} />*/}
                        {/*<Route path="users/:id" element={<UserDetailPage />} />*/}
                        {/*<Route path="jobs" element={<JobsListPage />} />*/}
                    </Route>

                    {/* Catch All */}
                    <Route path="*" element={<NotFoundPage />} />
                </Routes>
            </BrowserRouter>
        </ThemeProvider>
    );
}

export default App;