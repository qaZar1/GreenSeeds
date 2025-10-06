import React from "react";
import { List, Datagrid, TextField, EditButton, DeleteButton, SimpleList } from "react-admin";
import { useMediaQuery } from "@mui/material";
import { EmptyUser } from "./EmptyUser";
import UserListActions from "./Action";

const UserList = ({ ...props }) => {
    const isSmall = useMediaQuery((theme) => theme.breakpoints.down("sm"));

    return (
        <List
            resource="users"
            pagination={false}
            empty={<EmptyUser />}
            {...props}
            sx={{ padding: 2 }}
            actions={<UserListActions />}
        >
            {isSmall ? (
                <SimpleList
                    primaryText={record => `Пользователь: ${record.user}`}
                    secondaryText={record => (
                        <>
                            <span style={{ display: 'block' }}>Мин плотность: {record.min_density}</span>
                            <span style={{ display: 'block' }}>Макс плотность: {record.max_density}</span>
                            <span style={{ display: 'block' }}>Ёмкость бункера: {record.tank_capacity}</span>
                            <span style={{ display: 'block' }}>Задержка: {record.latency}</span>
                        </>
                    )}
                    
                    tertiaryText={record => (
                        <>
                            <EditButton record={record} />
                        </>
                    )}
                    rowClick={false}
                />
            ) : (
                <Datagrid rowClick="edit">
                    <TextField source="seed" label="Семена" />
                    <TextField source="min_density" label="Минимальная плотность" />
                    <TextField source="max_density" label="Максимальная плотность" />
                    <TextField source="tank_capacity" label="Емкость бункера" />
                    <TextField source="latency" label="Задержка" />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
};

export default UserList;