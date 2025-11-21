import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyReceipt } from "./EmptyReceipt";
import ReceiptListActions from "./Action";
import ReceiptListContent from "./Controller";

const ReceiptList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));

    return (
        <List
            resource="receipts"
            empty={<EmptyReceipt />}
            {...props}
            sx={{ padding: 2 }}
            actions={<ReceiptListActions />}
            title="Рецепты"
            pagination={false}
        >
            <ReceiptListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default ReceiptList;