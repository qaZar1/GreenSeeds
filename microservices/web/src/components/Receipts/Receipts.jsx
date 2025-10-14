import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyReceipt } from "./EmptyReceipt";
import ReceiptListActions from "./Action";
import ReceiptListContent from "./Controller";

const ReceiptList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    return (
        <List
            resource="receipts"
            empty={<EmptyReceipt />}
            {...props}
            sx={{ padding: 2 }}
            actions={<ReceiptListActions />}
        >
            <ReceiptListContent isSmall={isSmall} />
        </List>
    );
};

export default ReceiptList;