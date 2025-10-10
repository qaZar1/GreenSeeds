import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyBunker } from "./EmptyBunker";

const BunkerListContent = ({ isSmall }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyBunker />;

    return isSmall ? (
        <SimpleList
            primaryText={record => `Бункер: ${record.bunker}`}
            secondaryText={record => `Расстояние: ${record.distance}`}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptyBunker />}
        />
    ) : (
        <Datagrid
            rowClick="edit"
            bulkActionButtons={false}
            empty={<EmptyBunker />}
        >
            <TextField source="bunker" label="Бункер" />
            <TextField source="distance" label="Расстояние" />
            <EditButton />
        </Datagrid>
    )
};

export default BunkerListContent;
