import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyPlacement } from "./EmptyPlacement";
import PlacementListActions from "./Action";
import PlacementListContent from "./Controller";

const PlacementList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    return (
        <List
            resource="placements"
            pagination={false}
            empty={<EmptyPlacement />}
            {...props}
            sx={{ padding: 2 }}
            actions={<PlacementListActions />}
        >
            <PlacementListContent isSmall={isSmall} />
        </List>
    );
};

export default PlacementList;