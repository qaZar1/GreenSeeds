import React from "react";
import { List } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyUser } from "./EmptyUser";
import UserListActions from "./Action";
import { CreateButton } from "react-admin";
import { jwtDecode } from "jwt-decode";
import UserListContent from "./Controller";
import { useNotify } from "react-admin";

const UserList = (props) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    const notify = useNotify();

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
            title="Пользователи"
            queryOptions={{
                onError: () => notify("Ошибка загрузки пользователей", { type: "error" }),
            }}
        >
            <UserListContent isSmall={isSmall} isMedium={isMedium} currentUsername={currentUsername} />
        </List>
    );
};

export default UserList;
