// SeedCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, BooleanInput } from "react-admin";
import { ToolbarSave } from "../utils/Toolbars";
import BackButton from "../utils/Back";

const UserCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarSave />}>
                <TextInput source="username" label="Пользователь" />
                <TextInput source="full_name" label="Полное имя" />
                <BooleanInput source="is_admin" label="Администратор" />
            </SimpleForm>
        </Create>
    );
}

export default UserCreate;
