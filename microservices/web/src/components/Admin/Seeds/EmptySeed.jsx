import React from "react";
import { Box, Typography } from "@mui/material";
import SeedsListActions from "./Action";

export const EmptySeed = () => (
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
        <Typography variant="h6">Нет данных о семенах</Typography>
        <Typography variant="body2">Добавьте семена через панель администратора</Typography>
        <Box sx={{ mt: 2 }}>
            <SeedsListActions />
        </Box>
    </Box>
);
