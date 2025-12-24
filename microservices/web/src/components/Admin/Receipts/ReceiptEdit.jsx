// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";
import { useNotify } from "react-admin";

const ReceiptEdit = () => {
    const notify = useNotify();
    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Редактирование рецепта"
            mutationOptions={{
                onSuccess: () => notify("Рецепт изменён", { type: "success" }),
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />}>
                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed_ru" id="seed" label="Семена"/>
                </ReferenceInput>
                <TextInput source="gcode" label="Код запуска" multiline/>
                <TextInput source="description" label="Описание" />
            </SimpleForm>
        </Edit>
    );
};

export default ReceiptEdit;
