import React from "react";
import { useListContext } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { Box } from "@mui/material";
import Tasks from "./ActiveTasks";

const TasksListContent = () => {
    const { isLoading, ids, data } = useListContext();

    if (isLoading) return <LoadingOverlay />;

    return (
        <Box p={2} maxWidth={800} margin="auto">
            <Tasks tasks={data} />
        </Box>
    )
};

export default TasksListContent;
