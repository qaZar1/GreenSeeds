import React from "react";
import { useListContext, Datagrid, TextField, SimpleList, ShowButton, BooleanField } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EmptyReports } from "./EmptyReports";
import { DateField } from "react-admin";
import GroupedDatagrid from "./Grid";

const ReportsListContent = ({ isSmall, isMedium }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyReports />;

    return (
    //     <SimpleList
    //         primaryText={record => `Смена: ${record.shift}`}
    //         secondaryText={record => (
    //             <>
    //                 <span style={{ display: 'block' }}>Номер сменного задания: {record.number}</span>
    //                 <span style={{ display: 'block' }}>Рецепт: {record.receipt}</span>
    //                 <span style={{ display: 'block' }}>Дата: <DateField source="dt" showTime locales="ru-RU" /></span>
    //                 <span style={{ display: 'block' }}>Успешно: <BooleanField source="success" /></span>
    //             </>
    //         )}
    //         tertiaryText={record => (
    //             <>
    //                 <ShowButton record={record} label="Показать"/>
    //             </>
    //         )}
    //         rowClick={false}
    //         empty={<EmptyReports />}
    //     />
    // ) : (
        // <Datagrid
        //     rowClick={false}
        //     bulkActionButtons={false}
        //     empty={<EmptyReports />}
        // >
        //     <TextField source="shift" label="Смена" />
        //     <TextField source="number" label="Номер задания" />
        //     <TextField source="receipt" label="Рецепт" />
        //     <DateField source="dt" showTime locales="ru-RU" label="Дата"/>
        //     <BooleanField source="success" label="Успешно" />
        //     <ShowButton label="Показать"/>
        // </Datagrid>
        <GroupedDatagrid />
    )
};

export default ReportsListContent;
