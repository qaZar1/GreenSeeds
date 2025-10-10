import React from "react";
import { getToken } from "../../dataProvider";
import { Button } from "@mui/material";

// Кнопка сброса пароля
const ResetPasswordButton = ({ record, fullWidth = false }) => {
    const token = getToken();
    
    const resetPwd = async () => {
        const response = await fetch(`/api/users/change-password`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({ username: record.username }),
        });
    
        if (response.status !== 204) {
            throw new Error(`Ошибка сброса пароля для ${record.username}`);
        }
    }    

    return (
        <Button 
            variant="outlined" 
            color="primary" 
            size="small" 
            onClick={resetPwd}
            fullWidth={fullWidth}
        >
            Сбросить пароль
        </Button>
    );
};

export default ResetPasswordButton;