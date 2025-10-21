import React from "react";
import { Box, Typography } from "@mui/material";
import ShiftListActions from "./Action";

export const EmptyShift = () => (
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
        <Typography variant="h6">Нет данных о сменах</Typography>
        <Typography variant="body2">Добавьте смену через панель администратора</Typography>
        <Box sx={{ mt: 2 }}>
            <ShiftListActions />
        </Box>
    </Box>
);
