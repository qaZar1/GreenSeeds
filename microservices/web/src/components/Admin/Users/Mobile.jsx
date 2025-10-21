import React from "react";
import { useNotify, useRefresh } from "react-admin";
import { IconButton, Menu, MenuItem } from "@mui/material";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import { getToken } from "../../../dataProvider";
import ResetPasswordButton from "./ResetPwd";
import { Button } from "@mui/material";
import dataProvider from "../../../dataProvider";

const MobileActionsMenu = ({ record, currentUsername }) => {
    const [anchorEl, setAnchorEl] = React.useState(null);
    const provider = dataProvider;
    const notify = useNotify();
    const refresh = useRefresh();
    const open = Boolean(anchorEl);
    const token = getToken();
    

    const handleClick = (event) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleDelete = async () => {
        provider.delete("users", { id: record.username, previousData: record});
    }

    const handleToggle = async () => {
        const newValue = !record.is_admin;
        console.log(newValue) // инвертируем текущую роль

        try {
            await fetch(`/api/users/update`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    username: record.username,
                    is_admin: newValue,
                }),
            });
            notify("Роль обновлена", "success");
            refresh();
        } catch (err) {
            notify("Ошибка при обновлении роли", "error");
        }
    };

    return (
        <>
            <IconButton onClick={handleClick}>
                <MoreVertIcon />
            </IconButton>
            <Menu anchorEl={anchorEl} open={open} onClose={handleClose}>
                <MenuItem onClick={handleClose}>
                    <ResetPasswordButton record={record} fullWidth/>
                </MenuItem>

                <MenuItem
                    onClick={() => { handleClose(); handleDelete(); }}
                    disabled={record.username === currentUsername}
                >
                    <Button
                        variant="outlined"
                        color="primary"
                        size="small"
                        fullWidth
                    >
                        Удалить
                    </Button>
                </MenuItem>

                <MenuItem onClick={() => { 
                    if (record.username === currentUsername) return; // предотвращаем клик
                    handleClose(); 
                    handleToggle(); 
                    }}
                    disabled={record.username === currentUsername}
                >
                    <Button
                        variant="outlined"
                        color="primary"
                        size="small"
                        fullWidth
                    >
                        Поменять роль
                    </Button>
                </MenuItem>
            </Menu>
        </>
    );
};

export default MobileActionsMenu;