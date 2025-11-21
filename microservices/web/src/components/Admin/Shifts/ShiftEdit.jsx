import React from "react";
import { Edit, SimpleForm, DateTimeInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useParams } from "react-router-dom";

const ShiftEdit = () => {
    const { id } = useParams();
    console.log(id);
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic" title="Редактирование смены">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <DateTimeInput source="dt" label="Дата" />
            </SimpleForm>
        </Edit>
    );
};

export default ShiftEdit;