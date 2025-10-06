// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../utils/Toolbars";
import BackButton from "../utils/Back";

const BunkerEdit = () => {
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <NumberInput source="bunker" label="Бункер" disabled/>
                <NumberInput source="distance" label="Расстояние" />
            </SimpleForm>
        </Edit>
    );
};

export default BunkerEdit;
