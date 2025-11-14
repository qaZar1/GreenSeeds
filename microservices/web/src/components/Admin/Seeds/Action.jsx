import React from "react";
import { TopToolbar } from "react-admin";
import { CreateButton } from "react-admin";
import { alpha } from '@mui/material/styles';
import { useTheme } from '@mui/material/styles';

const SeedsListActions = () => {
    const theme = useTheme();
    return (
    <TopToolbar>
        <CreateButton
            label="ДОБАВИТЬ СЕМЕНА"
            sx={theme => ({
                textTransform: 'none',
                color: theme.palette.primary.main,
                '&:hover': {
                    bgcolor: theme.palette.mode === 'light'
                        ? alpha(theme.palette.primary.main, 0.15)
                        : alpha(theme.palette.primary.main, 0.1),
                    color: theme.palette.primary.main,
                },
            })}
        />
    </TopToolbar>
    );
};

export default SeedsListActions;