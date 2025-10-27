import React from "react";
import { Box, Typography } from "@mui/material";
import BunkerListActions from "./Action";

export const EmptyBunker = () => (
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
        <Typography variant="h6">Нет данных о бункерах</Typography>
        <Typography variant="body2">Добавьте бункеры через панель администратора</Typography>
        <Box sx={{ mt: 2 }}>
            <BunkerListActions />
        </Box>
    </Box>
);
