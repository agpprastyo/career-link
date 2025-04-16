import { useState } from 'react';
import { Outlet, useNavigate, Link } from 'react-router-dom';
import {
    Box,
    Drawer,
    AppBar,
    Toolbar,
    List,
    Typography,
    Divider,
    IconButton,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Avatar,
    Menu,
    MenuItem,
    Tooltip
} from '@mui/material';
import {
    Menu as MenuIcon,
    Dashboard as DashboardIcon,
    People as PeopleIcon,
    Work as WorkIcon,
    Business as BusinessIcon,
    AccountCircle as AccountCircleIcon,
    Logout as LogoutIcon,
    Settings as SettingsIcon,
    ChevronLeft as ChevronLeftIcon
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { logout } from '../../store/slices/authSlice';

const drawerWidth = 240;

const AdminLayout = () => {
    const [open, setOpen] = useState(true);
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const navigate = useNavigate();
    const dispatch = useAppDispatch();
    const { user } = useAppSelector(state => state.auth);

    const handleDrawerToggle = () => {
        setOpen(!open);
    };

    const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleProfileMenuClose = () => {
        setAnchorEl(null);
    };

    const handleLogout = () => {
        dispatch(logout());
        navigate('/login');
    };

    const menuItems = [
        { text: 'Dashboard', icon: <DashboardIcon />, path: '/dashboard' },
        { text: 'Users', icon: <PeopleIcon />, path: '/users' },
        { text: 'Jobs', icon: <WorkIcon />, path: '/jobs' },
        { text: 'Companies', icon: <BusinessIcon />, path: '/companies' },
    ];

    return (
        <Box sx={{ display: 'flex' }}>
            <AppBar position="fixed"
                    sx={{
                        zIndex: theme => theme.zIndex.drawer + 1,
                        transition: theme => theme.transitions.create(['width', 'margin'], {
                            easing: theme.transitions.easing.sharp,
                            duration: theme.transitions.duration.leavingScreen,
                        }),
                        ...(open && {
                            marginLeft: drawerWidth,
                            width: `calc(100% - ${drawerWidth}px)`,
                            transition: theme => theme.transitions.create(['width', 'margin'], {
                                easing: theme.transitions.easing.sharp,
                                duration: theme.transitions.duration.enteringScreen,
                            }),
                        }),
                    }}
            >
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={handleDrawerToggle}
                        sx={{ mr: 2 }}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
                        Career Link Admin
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Tooltip title="Account settings">
                            <IconButton
                                size="large"
                                edge="end"
                                aria-label="account of current user"
                                aria-controls="menu-appbar"
                                aria-haspopup="true"
                                onClick={handleProfileMenuOpen}
                                color="inherit"
                            >
                                {user?.avatar ? (
                                    <Avatar src={user.avatar} alt={user?.username || 'User'} />
                                ) : (
                                    <AccountCircleIcon />
                                )}
                            </IconButton>
                        </Tooltip>
                    </Box>
                </Toolbar>
            </AppBar>

            <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                }}
                open={Boolean(anchorEl)}
                onClose={handleProfileMenuClose}
            >
                <MenuItem onClick={() => navigate('/profile')}>
                    <ListItemIcon><AccountCircleIcon fontSize="small" /></ListItemIcon>
                    <ListItemText>Profile</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => navigate('/settings')}>
                    <ListItemIcon><SettingsIcon fontSize="small" /></ListItemIcon>
                    <ListItemText>Settings</ListItemText>
                </MenuItem>
                <Divider />
                <MenuItem onClick={handleLogout}>
                    <ListItemIcon><LogoutIcon fontSize="small" /></ListItemIcon>
                    <ListItemText>Logout</ListItemText>
                </MenuItem>
            </Menu>

            <Drawer
                variant="permanent"
                open={open}
                sx={{
                    width: drawerWidth,
                    flexShrink: 0,
                    '& .MuiDrawer-paper': {
                        width: drawerWidth,
                        boxSizing: 'border-box',
                        ...(open ? {
                            width: drawerWidth,
                            transition: theme => theme.transitions.create('width', {
                                easing: theme.transitions.easing.sharp,
                                duration: theme.transitions.duration.enteringScreen,
                            }),
                            overflowX: 'hidden',
                        } : {
                            transition: theme => theme.transitions.create('width', {
                                easing: theme.transitions.easing.sharp,
                                duration: theme.transitions.duration.leavingScreen,
                            }),
                            overflowX: 'hidden',
                            width: theme => theme.spacing(7),
                            [theme.breakpoints.up('sm')]: {
                                width: theme => theme.spacing(9),
                            },
                        }),
                    },
                }}
            >
                <Toolbar
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'flex-end',
                        px: [1],
                    }}
                >
                    <IconButton onClick={handleDrawerToggle}>
                        <ChevronLeftIcon />
                    </IconButton>
                </Toolbar>
                <Divider />
                <List component="nav">
                    {menuItems.map((item) => (
                        <ListItem key={item.text} disablePadding sx={{ display: 'block' }}>
                            <ListItemButton
                                component={Link}
                                to={item.path}
                                sx={{
                                    minHeight: 48,
                                    justifyContent: open ? 'initial' : 'center',
                                    px: 2.5,
                                }}
                            >
                                <ListItemIcon
                                    sx={{
                                        minWidth: 0,
                                        mr: open ? 3 : 'auto',
                                        justifyContent: 'center',
                                    }}
                                >
                                    {item.icon}
                                </ListItemIcon>
                                <ListItemText
                                    primary={item.text}
                                    sx={{ opacity: open ? 1 : 0 }}
                                />
                            </ListItemButton>
                        </ListItem>
                    ))}
                </List>
            </Drawer>

            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    p: 3,
                    width: { sm: `calc(100% - ${drawerWidth}px)` },
                    mt: '64px',
                }}
            >
                <Outlet />
            </Box>
        </Box>
    );
};

export default AdminLayout;