import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyShift } from "./EmptyShift";
import ShiftListActions from "./Action";
import ShiftListContent from "./Controller";
import { useNotify } from "react-admin";

const ShiftList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    const notify = useNotify();

    return (
        <List
            resource="shifts"
            empty={<EmptyShift />}
            {...props}
            sx={{ padding: 2 }}
            actions={<ShiftListActions />}
            title="План производства"
            pagination={false}
            queryOptions={{
                onError: () => notify("Ошибка загрузки смен", { type: "error" }),
            }}
        >
            <ShiftListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default ShiftList;