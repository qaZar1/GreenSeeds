// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, AutocompleteInput } from "react-admin";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const ReceiptCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание рецепта"
            mutationOptions={{
                onSuccess: (response) => {
                    notify("Рецепт создан", { type: "success" })
                    redirect('edit', 'receipts', response.id);
                },
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed_ru" id="seed" label="Семена"/>
                </ReferenceInput>
                <TextInput source="gcode" label="Код запуска" multiline/>
                <TextInput source="description" label="Описание" />
            </SimpleForm>
        </Create>
    );
}

export default ReceiptCreate;
