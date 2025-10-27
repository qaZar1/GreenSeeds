import React from "react";
import { List } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyReports } from "./EmptyReports";
import ReportsListContent from "./Controller";

const ReportsList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));

    return (
        <List
            resource="reports"
            empty={<EmptyReports />}
            {...props}
            pagination={false}
            perPage={false}
            sx={{ padding: 2 }}
            actions={false}
        >
            <ReportsListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default ReportsList;