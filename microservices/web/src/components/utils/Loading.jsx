// src/components/LoadingOverlay.jsx
import React from "react";
import { Box, CircularProgress, Typography } from "@mui/material";

export const LoadingOverlay = () => (
    <Box
        sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            height: "100%",
            minHeight: "300px",
            textAlign: "center",
            gap: 2,
        }}
    >
        <CircularProgress />
        <Typography variant="body1">Загрузка данных...</Typography>
    </Box>
);
