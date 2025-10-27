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
        <GroupedDatagrid />
    )
};

export default ReportsListContent;
