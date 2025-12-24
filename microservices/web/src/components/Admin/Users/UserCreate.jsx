import React from "react";
import { Create, SimpleForm, TextInput, BooleanInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const UserCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание пользователя"
            mutationOptions={{
                onSuccess: () => {
                    notify("Пользователь создан", { type: "success" });
                    redirect('list', 'users');
                },
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <TextInput source="username" label="Пользователь" />
                <TextInput source="full_name" label="Полное имя" />
                <BooleanInput source="is_admin" label="Администратор" />
            </SimpleForm>
        </Create>
    );
}

export default UserCreate;
