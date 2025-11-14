import React from "react";
import { useNotify, useRefresh } from "react-admin";
import { Switch } from "@mui/material";
import { getToken } from "../../../dataProvider";

const AdminToggleField = ({ record, currentUsername }) => {
    const notify = useNotify();
    const refresh = useRefresh();

    const handleChange = async (e) => {
        const newValue = e.target.checked;
        const token = getToken();

        try {
            await fetch(`/api/users/update`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({
                    username: record.username,
                    full_name: record.full_name,
                    is_admin: newValue,
                }),
            });
            notify('Роль обновлена', { type: 'success' });
            refresh();
        } catch (err) {
            notify('Ошибка при обновлении роли', { type: 'error' });
        }
    };

    return <Switch
        checked={record.is_admin}
        onChange={handleChange}
        color="primary"
        disabled={record.username === currentUsername}/>;
};

export default AdminToggleField;