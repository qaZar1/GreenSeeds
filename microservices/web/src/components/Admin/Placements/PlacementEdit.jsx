// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";

const PlacementEdit = () => {
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <ReferenceInput source="bunker" reference="bunkers">
                    <AutocompleteInput optionText="bunker" id="bunker" label="Бункер"/>
                </ReferenceInput>

                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed" id="seed" label="Семена"/>
                </ReferenceInput>
            </SimpleForm>
        </Edit>
    );
};

export default PlacementEdit;
