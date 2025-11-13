import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyBunker } from "./EmptyBunker";
import BunkerListActions from "./Action";
import BunkerListContent from "./Controller";

const BunkerList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    return (
        <List
            resource="bunkers"
            pagination={false}
            empty={<EmptyBunker />}
            {...props}
            sx={{ padding: 2 }}
            actions={<BunkerListActions />}
            title="Бункеры"
        >
            <BunkerListContent isSmall={isSmall} />
        </List>
    );
};

export default BunkerList;