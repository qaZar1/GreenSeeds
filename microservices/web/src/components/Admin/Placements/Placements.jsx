import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyPlacement } from "./EmptyPlacement";
import PlacementListActions from "./Action";
import PlacementListContent from "./Controller";

const PlacementList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));

    return (
        <List
            resource="placements"
            pagination={false}
            empty={<EmptyPlacement />}
            {...props}
            sx={{ padding: 2 }}
            actions={<PlacementListActions />}
            title="Расположение семян"
        >
            <PlacementListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default PlacementList;