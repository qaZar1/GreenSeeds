import React from "react";
import { List } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptySettings } from "./EmptySettings";
import SettingsListActions from "./Action";
import SettingsListContent from "./Controller"
import { useNotify } from "react-admin";

const SettingsList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));
    const isMedium = useMediaQuery((theme) => theme.breakpoints.between("sm", "md"));
    const notify = useNotify();

    return (
        <List
            resource="device-settings"
            pagination={false}
            empty={<EmptySettings />}
            {...props}
            sx={{ padding: 2 }}
            actions={<SettingsListActions />}
            title="Настройки устройств"
            queryOptions={{
                onError: () => notify("Ошибка загрузки настроек", { type: "error" })
            }}
        >
            <SettingsListContent isSmall={isSmall} isMedium={isMedium} />
        </List>
    );
};

export default SettingsList;