import React from "react";
import { Box, Typography } from "@mui/material";
import SettingsListActions from "./Action";

export const EmptySettings = () => (
    <Box
        sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            height: "200px",
            textAlign: "center",
            color: "text.secondary",
        }}
    >
        <Typography variant="h6">Нет данных о настройках</Typography>
        <Typography variant="body2">Добавьте настройки через панель администратора</Typography>
        <Box sx={{ mt: 2 }}>
            <SettingsListActions />
        </Box>
    </Box>
);
