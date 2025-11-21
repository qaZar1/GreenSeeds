import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptySeed } from "./EmptySeed";
import SeedListActions from "./Action";
import SeedListContent from "./Controller";

const SeedList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));

    return (
        <List
            resource="seeds"
            empty={<EmptySeed />}
            {...props}
            sx={{ padding: 2 }}
            actions={<SeedListActions />}
            title="Семена"
            pagination={false}
        >
            <SeedListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default SeedList;