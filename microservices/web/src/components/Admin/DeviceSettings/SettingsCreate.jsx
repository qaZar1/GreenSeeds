import React from "react";
import { Create, SimpleForm, TextInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";

const SettingsCreate = () => {
    const notify = useNotify();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание настроек устройства"
            mutationOptions={{
                onSuccess: () => notify("Настройки созданы", { type: "success" }),
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <TextInput source="key" label="Ключ"/>
                <TextInput source="value" label="Значение" multiline={true}/>
            </SimpleForm>
        </Create>
    );
}

export default SettingsCreate;
