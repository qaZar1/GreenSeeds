import React from "react";
import { Box, Typography } from "@mui/material";
import PlacementListActions from "./Action";

export const EmptyPlacement = () => (
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
        <Typography variant="h6">Нет данных о расположении семян в бункерах</Typography>
        <Typography variant="body2">Добавьте информацию о семенах в бункерах через панель администратора</Typography>
        <Box sx={{ mt: 2 }}>
            <PlacementListActions />
        </Box>
    </Box>
);
