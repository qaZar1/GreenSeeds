import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyReceipt } from "./EmptyReceipt";
import { DateField } from "react-admin";

const ReceiptListContent = ({ isSmall }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyReceipt />;

    return isSmall ? (
        <SimpleList
            primaryText={record => `Название: ${record.receipt}`}
            secondaryText={record => (
                <>
                    <span style={{ display: 'block' }}>Семена: {record.seed}</span>
                    <span>Обновлено: <DateField source="updated" showTime locales="ru-RU" /></span>
                    <span style={{ display: 'block' }}>Описание: {record.description}</span>
                </>
            )}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptyReceipt />}
        />
    ) : (
        <Datagrid
            rowClick="edit"
            bulkActionButtons={false}
            empty={<EmptyReceipt />}
        >
            <TextField source="receipt" label="Название" />
            <TextField source="seed" label="Семена" />
            <DateField source="updated" showTime locales="ru-RU" label="Обновлено"/>
            <TextField source="description" label="Описание" />
            {/* <EditButton /> */}
        </Datagrid>
    )
};

export default ReceiptListContent;
