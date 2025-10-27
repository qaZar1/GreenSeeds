import React from "react";
import { Box, Typography } from "@mui/material";

const EmptyTask = () => {
    return (
        <Box
            display="flex"
            justifyContent="center"
            alignItems="center"
            minHeight="50vh"
        >
            <Typography variant="h6" color="text.secondary">
                Сменных заданий нет
            </Typography>
        </Box>
    );
}

export default EmptyTask;