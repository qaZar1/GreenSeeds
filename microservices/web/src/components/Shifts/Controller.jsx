import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyShift } from "./EmptyShift";
import { DateField } from "react-admin";

const ShiftListContent = ({ isSmall }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyShift />;

    return isSmall ? (
        <SimpleList
            primaryText={<DateField source="dt" label="Дата" showTime={false}/>}
            secondaryText={record => `Оператор: ${record.username ? record.username : "Не определено"}`}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptyShift />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptyShift />}
        >
            <DateField source="dt" label="Дата" showTime={false}/>
            <TextField source="username" label="Пользователь" />
            <EditButton label="Редактировать" />
        </Datagrid>
    )
};

export default ShiftListContent;
