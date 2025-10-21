import React from "react";
import { TopToolbar } from "react-admin";
import { CreateButton } from "react-admin";
import { alpha } from '@mui/material/styles';
import { useTheme } from '@mui/material/styles';

const AssignmentsListActions = () => {
    const theme = useTheme();
    return (
    <TopToolbar>
        <CreateButton
            label="ДОБАВИТЬ СМЕННОЕ ЗАДАНИЕ"
            sx={theme => ({
                textTransform: 'none',
                color: theme.palette.primary.main,
                '&:hover': {
                    bgcolor: theme.palette.mode === 'light'
                        ? alpha(theme.palette.primary.main, 0.15) // чуть ярче на светлой теме
                        : alpha(theme.palette.primary.main, 0.1), // чуть мягче на темной
                    color: theme.palette.primary.main,
                },
            })}
        />
    </TopToolbar>
    );
};

export default AssignmentsListActions;