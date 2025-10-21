// ShiftEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput, DateTimeInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, AutocompleteInput } from "react-admin";

const ShiftEdit = () => {
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <DateTimeInput source="dt" label="Дата"/>
            </SimpleForm>
        </Edit>
    );
};

export default ShiftEdit;
