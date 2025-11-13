import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyShift } from "./EmptyShift";
import ShiftListActions from "./Action";
import ShiftListContent from "./Controller";

const ShiftList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    return (
        <List
            resource="shifts"
            empty={<EmptyShift />}
            {...props}
            sx={{ padding: 2 }}
            actions={<ShiftListActions />}
            title="План производства"
            pagination={false}
        >
            <ShiftListContent isSmall={isSmall} />
        </List>
    );
};

export default ShiftList;