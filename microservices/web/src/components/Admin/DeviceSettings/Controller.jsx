import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EditButton } from "react-admin";
import { EmptySettings } from "./EmptySettings";

const SettingsListContent = ({ isSmall, isMedium }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptySettings />;

    return isSmall || isMedium ? (
        <SimpleList
            primaryText={record => `Ключ: ${record.key}`}
            secondaryText={record => `Значение: ${record.value}`}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptySettings />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptySettings />}
        >
            <TextField source="key" label="Ключ" />
            <TextField source="value" label="Значение" />
            <EditButton label="Редактировать" />
        </Datagrid>
    )
};

export default SettingsListContent;
