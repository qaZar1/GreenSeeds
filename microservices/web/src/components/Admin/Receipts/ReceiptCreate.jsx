// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, AutocompleteInput } from "react-admin";

const ReceiptCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarSave />}>
                <ReferenceInput source="seed_ru" reference="seeds">
                    <AutocompleteInput optionText="seed_ru" id="seed" label="Семена"/>
                </ReferenceInput>
                <TextInput source="gcode" label="Код запуска" />
                <TextInput source="description" label="Описание" />
            </SimpleForm>
        </Create>
    );
}

export default ReceiptCreate;
