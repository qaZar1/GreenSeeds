import React from "react";
import { Box, Typography } from "@mui/material";
import { CreateButton } from "react-admin";
import ReceiptListActions from "./Action";

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
        <Typography variant="h6">Нет данных о рецептах</Typography>
        <Typography variant="body2">Добавьте рецепты через панель администратора</Typography>
        <Box sx={{ mt: 2 }}>
            <ReceiptListActions />
        </Box>
    </Box>
);
