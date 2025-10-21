import React from "react";
import { Create, SimpleForm, DateTimeInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";

const ShiftCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarSave />}>
                <DateTimeInput source="dt" label="Дата"/>
            </SimpleForm>
        </Create>
    );
}

export default ShiftCreate;
