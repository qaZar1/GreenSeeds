import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList, useNotify } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyAssignments } from "./EmptyAssign";
import AssignmentsListActions from "./Action";
import AssignmentsListContent from "./Controller";

const AssignmentsList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    const notify = useNotify();

    return (
        <List
            resource="assignments"
            empty={<EmptyAssignments />}
            {...props}
            sx={{ padding: 2 }}
            actions={<AssignmentsListActions />}
            title="Сменные задания"
            pagination={false}
            queryOptions={{
                onError: () => notify("Ошибка загрузки заданий", { type: "error" })
            }}
        >
            <AssignmentsListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default AssignmentsList;