// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const BunkerCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание бункера"
            mutationOptions={{
                onSuccess: (response) => {
                    notify("Бункер создан", { type: "success" })
                    redirect('edit', 'bunkers', response.id);
                },
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <NumberInput source="bunker" label="Бункер" max={10} min={1}/>
                <NumberInput source="distance" label="Расстояние" />
            </SimpleForm>
        </Create>
    );
}

export default BunkerCreate;
