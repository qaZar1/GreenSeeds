import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../utils/Loading";
import { EditButton } from "react-admin";
import { EmptySeed } from "./EmptySeed";

const SeedListContent = ({ isSmall, isMedium }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay/>;
    if (error) return <EmptySeed />;

    return isSmall || isMedium ? (
        <SimpleList
            primaryText={record => `Семена: ${record.seed}`}
            secondaryText={record => (
                <>
                    <span style={{ display: 'block' }}>Мин плотность: {record.min_density}</span>
                    <span style={{ display: 'block' }}>Макс плотность: {record.max_density}</span>
                    <span style={{ display: 'block' }}>Количество семян в бункере: {record.tank_capacity}</span>
                    <span style={{ display: 'block' }}>Задержка: {record.latency}</span>
                </>
            )}
            
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptySeed />}
        />
    ) : (
        <Datagrid
            rowClick="edit"
            bulkActionButtons={false}
            empty={<EmptySeed />}
        >
            <TextField source="seed" label="Семена" />
            <TextField source="min_density" label="Минимальная плотность" />
            <TextField source="max_density" label="Максимальная плотность" />
            <TextField source="tank_capacity" label="Количество семян в бункере" />
            <TextField source="latency" label="Задержка" />
            <EditButton />
        </Datagrid>
    )
};

export default SeedListContent;
