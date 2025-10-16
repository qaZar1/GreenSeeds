import React from "react";
import { Box, Typography } from "@mui/material";

export const EmptyReports = () => (
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
        <Typography variant="h6">Нет данных</Typography>
    </Box>
);
