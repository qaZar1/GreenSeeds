// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, BooleanInput } from "react-admin";
import { ToolbarEdit } from "../utils/Toolbars";
import BackButton from "../utils/Back";

const UserEdit = () => {
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <TextInput source="username" label="Пользователь" disabled/>
                <TextInput source="full_name" label="Полное имя" />
                <BooleanInput source="is_admin" label="Администратор" />
            </SimpleForm>
        </Edit>
    );
};

export default UserEdit;
