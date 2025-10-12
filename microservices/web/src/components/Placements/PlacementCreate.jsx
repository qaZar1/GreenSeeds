// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../utils/Toolbars";
import BackButton from "../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";

const PlacementCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarSave />}>
                <ReferenceInput source="bunker" reference="bunkers">
                    <AutocompleteInput optionText="bunker" id="bunker" label="Бункер"/>
                </ReferenceInput>

                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed" id="seed" label="Семена"/>
                </ReferenceInput>
            </SimpleForm>
        </Create>
    );
}

export default PlacementCreate;
