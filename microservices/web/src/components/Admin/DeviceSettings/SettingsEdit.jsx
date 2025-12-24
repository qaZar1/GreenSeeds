import React from "react";
import { Edit, SimpleForm, TextInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";

const SettingsEdit = () => {
    const notify = useNotify();

    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Редактирование настроек"
            mutationOptions={{
                onSuccess: () => notify("Настройки изменены", { type: "success" }),
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />}>
                <TextInput source="key" label="Ключ" disabled/>
                <TextInput source="value" label="Значение" multiline={true}/>
            </SimpleForm>
        </Edit>
    );
};

export default SettingsEdit;
