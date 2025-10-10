import React from "react";
import { 
    List, 
    Datagrid, 
    TextField, 
    DeleteButton,
    FunctionField
} from "react-admin";
import { useMediaQuery, Button } from "@mui/material";
import { EmptyUser } from "./EmptyUser";
import UserListActions from "./Action";
import { CreateButton } from "react-admin";
import { jwtDecode } from "jwt-decode";
import { IconButton, Menu, MenuItem } from "@mui/material";
import { useDataProvider } from "react-admin";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import { useNotify } from "react-admin";
import { useRefresh } from "react-admin";
import { getToken } from "../../dataProvider";
import UserListContent from "./Controller";

const UserList = (props) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    let currentUsername = null;
    try {
        const stored = localStorage.getItem("auth");
        if (stored) {
            const parsed = JSON.parse(stored);
            if (parsed?.token) {
                const decoded = jwtDecode(parsed.token);
                currentUsername = decoded?.username || null;
            }
        }
    } catch (e) {
        console.warn("Ошибка получения username из localStorage:", e);
    }

    return (
        <List
            resource="users"
            pagination={false}
            empty={<EmptyUser />}
            {...props}
            sx={{ padding: 2 }}
            actions={isSmall ? <CreateButton /> : <UserListActions />}
        >
            <UserListContent isSmall={isSmall} currentUsername={currentUsername} />
        </List>
    );
};

export default UserList;
