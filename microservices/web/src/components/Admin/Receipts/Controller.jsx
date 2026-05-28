import React from "react";
import { useListContext, Datagrid, TextField, SimpleList } from "react-admin";
import { LoadingOverlay } from "../../utils/Loading";
import { EditButton } from "react-admin";
import { EmptyRecipe } from "./EmptyRecipe";
import { DateField } from "react-admin";

const RecipeListContent = ({ isSmall, isMedium }) => {
    const { isLoading, ids, data, error } = useListContext();

    if (isLoading) return <LoadingOverlay />;
    if (error) return <EmptyRecipe />;

    return isSmall || isMedium ? (
        <SimpleList
            primaryText={record => `Семена: ${record.seed_ru}`}
            secondaryText={record => (
                <>
                    <span style={{ display: 'block' }}>Описание: {record.description}</span>
                    <span style={{ display: 'block' }}>Обновлено: <DateField
                        source="updated"
                        showTime
                        locales="ru-RU"
                    /></span>
                </>
            )}
            tertiaryText={record => (
                <>
                    <EditButton record={record} />
                </>
            )}
            rowClick={false}
            empty={<EmptyRecipe />}
        />
    ) : (
        <Datagrid
            rowClick={false}
            bulkActionButtons={false}
            empty={<EmptyRecipe />}
        >
            <TextField source="seed_ru" label="Семена" />
            <TextField source="description" label="Описание" />
            <DateField source="updated" showTime locales="ru-RU" label="Обновлено"/>
            <EditButton label="Редактировать"/>
        </Datagrid>
    )
};

export default RecipeListContent;
