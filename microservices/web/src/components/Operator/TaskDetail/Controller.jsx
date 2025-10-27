import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyDetail } from "./EmptyDetail";

const BunkerListContent = ({ isSmall }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyDetail />;

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
            empty={<EmptyDetail />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptyDetail />}
        >
            <TextField source="bunker" label="Бункер" />
            <TextField source="distance" label="Расстояние" />
            <EditButton label="Редактировать" />
        </Datagrid>
    )
};

export default BunkerListContent;
