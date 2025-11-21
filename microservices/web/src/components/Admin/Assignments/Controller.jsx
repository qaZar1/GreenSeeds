import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyAssignments } from "./EmptyAssign";
import { DateField } from "react-admin";


const AssignmentsListContent = ({ isSmall, isMedium }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyAssignments />;

    return isSmall || isMedium ? (
        <SimpleList
            primaryText={record => `Смена: ${record.shift}`}
            secondaryText={record => (
                <>
                    <span style={{ display: 'block' }}>Номер сменного задания: {record.number}</span>
                    <span style={{ display: 'block' }}>Название рецепта: {record.description}</span>
                    <span style={{ display: 'block' }}>Количество (лотков): {record.amount}</span>
                    
                </>
            )}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptyAssignments />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptyAssignments />}
        >
            <TextField source="shift" label="Номер смены" />
            <TextField source="number" label="Номер сменного задания" />
            <TextField source="description" label="Название рецепта" />
            <TextField source="amount" label="Количество (лотков)"/>
            <EditButton label="Редактировать"/>
        </Datagrid>
    )
};

export default AssignmentsListContent;
