import React from "react";
import { Box, Typography } from "@mui/material";

export const EmptyReceipt = () => (
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
        <Typography variant="h6">Нет данных о семенах в бункере</Typography>
        <Typography variant="body2">Добавьте семена в бункер через панель администратора</Typography>
    </Box>
);
