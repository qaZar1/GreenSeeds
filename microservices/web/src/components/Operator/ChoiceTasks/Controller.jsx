import React from "react";
import { useListContext } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import Task from "./FreeTask";
import { Box } from "@mui/material";
import EmptyChoice from "./EmptyChoice";

const ChoiceListContent = ({ username }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyChoice />;

    return (
        <Box
            display="grid"
            gridTemplateColumns="repeat(auto-fit, minmax(350px, 1fr))"
            gap={2}
            p={2}
            >
            {data.map((task) => (
                <Task key={task.id} task={task} username={username} />
            ))}
        </Box>
    )
};

export default ChoiceListContent;
