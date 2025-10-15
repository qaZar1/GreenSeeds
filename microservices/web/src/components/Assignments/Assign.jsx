import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyAssignments } from "./EmptyAssign";
import AssignmentsListActions from "./Action";
import AssignmentsListContent from "./Controller";

const AssignmentsList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    return (
        <List
            resource="assignments"
            empty={<EmptyAssignments />}
            {...props}
            sx={{ padding: 2 }}
            actions={<AssignmentsListActions />}
        >
            <AssignmentsListContent isSmall={isSmall} />
        </List>
    );
};

export default AssignmentsList;