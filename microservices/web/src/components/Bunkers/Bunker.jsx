import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyBunker } from "./EmptyBunker";
import BunkerListActions from "./Action";

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
        >
            {isSmall ? (
                <SimpleList
                    primaryText={record => `Бункер: ${record.bunker}`}
                    secondaryText={record => `Расстояние: ${record.distance}`}
                    tertiaryText={record => (
                        <>
                            <EditButton record={record} />
                        </>
                    )}
                    rowClick={false}
                />
            ) : (
                <Datagrid rowClick="edit">
                    <TextField source="bunker" label="Бункер" />
                    <TextField source="distance" label="Расстояние" />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
};

export default BunkerList;