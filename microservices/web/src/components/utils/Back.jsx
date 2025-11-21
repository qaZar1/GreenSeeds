import React from "react";
import { useNavigate } from "react-router-dom";
import { IconButton } from "@mui/material";
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { TopToolbar } from "react-admin";
import { useTheme } from "@mui/material/styles";
import { useMediaQuery } from "@mui/material";
import { useResourceContext } from "react-admin";
import { Box } from "@mui/material";

const BackButton = ({ to }) => {
    const navigate = useNavigate();
    const theme = useTheme();
    const isSmall = useMediaQuery(theme.breakpoints.down('sm'));
    const resource = useResourceContext();

    const handleClick = () => {
        navigate(to ? to : `/${resource}`);
    };

    return (
        <IconButton
            variant="contained"
            color="primary"
            onClick={handleClick}
            sx={{
                minWidth: isSmall ? 36 : 40,
                width: isSmall ? 36 : 40,
                height: isSmall ? 36 : 40,
                borderRadius: '50%',
                padding: 0,
                mb: 2,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                '& .MuiButton-startIcon': {
                    margin: 0,
                },
            }}
        >
            <ArrowBackIcon fontSize={isSmall ? 'small' : 'medium'} />
        </IconButton>
    );
};

const TopToolbarWithBackButton = ({ to }) => {
    return (
        <TopToolbar
            sx={{ justifyContent: 'flex-start', p: 0 }}
            component={Box}
        >
            <BackButton to={to} />
        </TopToolbar>
    );
};

export default TopToolbarWithBackButton;