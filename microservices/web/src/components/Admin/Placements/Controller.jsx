import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyPlacement } from "./EmptyPlacement";

const PlacementListContent = ({ isSmall }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyPlacement />;

    return isSmall ? (
        <SimpleList
            primaryText={record => `Бункер: ${record.bunker}`}
            secondaryText={record => `Семена: ${record.seed_ru}`}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptyPlacement />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptyPlacement />}
        >
            <TextField source="bunker" label="Бункер" />
            <TextField source="seed_ru" label="Семена" />
            <EditButton label="Редактировать" />
        </Datagrid>
    )
};

export default PlacementListContent;
