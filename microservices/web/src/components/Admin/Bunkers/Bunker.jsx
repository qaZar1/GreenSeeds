import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyBunker } from "./EmptyBunker";
import BunkerListActions from "./Action";
import BunkerListContent from "./Controller";
import { useNotify } from "react-admin";

const BunkerList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    const notify = useNotify();

    return (
        <List
            resource="bunkers"
            pagination={false}
            empty={<EmptyBunker />}
            {...props}
            sx={{ padding: 2 }}
            actions={<BunkerListActions />}
            title="Бункеры"
            queryOptions={{
                onError: () => notify("Ошибка загрузки бункеров", { type: "error" })
            }}
        >
            <BunkerListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default BunkerList;