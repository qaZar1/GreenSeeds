import React from "react";
import { Toolbar, SaveButton, DeleteButton } from "react-admin";
import { Box } from "@mui/material";

const ToolbarSave = (props) => (
    <Toolbar {...props}>
        <SaveButton label="Сохранить" />
    </Toolbar>
);

const ToolbarEdit = (props) => (
    <Toolbar {...props}>
        <Box
            sx={{
                display: 'flex',
                justifyContent: 'space-between',
                width: '100%',
            }}
        >
            <SaveButton label="Сохранить" />
            <DeleteButton
                label="Удалить"
                confirmTitle="Удалить элемент?"
                confirmContent="Это действие нельзя отменить. Точно удалить?"
            />
        </Box>
    </Toolbar>
);

export { ToolbarSave, ToolbarEdit };
