// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";
import { useNotify } from "react-admin";

const PlacementEdit = () => {
    const notify = useNotify();
    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Редактирование расположения семян"
            mutationOptions={{
                onSuccess: () => notify("Расположение семян изменено", { type: "success" }),
                onError: () => notify("Невозможно обновить расположение", { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />}>
                <ReferenceInput source="bunker" reference="bunkers">
                    <AutocompleteInput optionText="bunker" id="bunker" label="Бункер"/>
                </ReferenceInput>

                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed_ru" id="seed" label="Семена"/>
                </ReferenceInput>
                <NumberInput source="amount" label="Количество семян в бункере(шт лотков)" min={0} defaultValue={0}/>
            </SimpleForm>
        </Edit>
    );
};

export default PlacementEdit;
