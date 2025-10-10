import React from "react";
import { AppBar, Layout, UserMenu, Logout } from "react-admin";
import { MenuItem } from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import { useNavigate } from "react-router-dom";

const CustomUserMenu = (props) => {
    const navigate = useNavigate();

    return(
    <UserMenu {...props}>
        <MenuItem onClick={() => navigate('/profile')}>
            <SettingsIcon style={{ marginRight: 8 }} />
            Профиль
        </MenuItem>
        <Logout />
    </UserMenu>
    );
};

const CustomAppBar = (props) => (
    <AppBar {...props} userMenu={<CustomUserMenu />} />
);

const CustomLayout = (props) => (
    <Layout {...props} appBar={CustomAppBar} />
);

export default CustomLayout;
