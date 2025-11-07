import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EditButton } from "react-admin";
import { EmptySeed } from "./EmptySeed";

const SeedListContent = ({ isSmall, isMedium }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay/>;
    if (error) return <EmptySeed />;

    return isSmall || isMedium ? (
        <SimpleList
            primaryText={record => `Семена: ${record.seed_ru}`}
            secondaryText={record => (
                <>
                    <span style={{ display: 'block' }}>Мин плотность: {record.min_density}</span>
                    <span style={{ display: 'block' }}>Количество семян в бункере: {record.tank_capacity}</span>
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
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptySeed />}
        >
            <TextField source="seed_ru" label="Семена" />
            <TextField source="tank_capacity" label="Количество семян в бункере" />
            <EditButton label="Редактировать" />
        </Datagrid>
    )
};

export default SeedListContent;
