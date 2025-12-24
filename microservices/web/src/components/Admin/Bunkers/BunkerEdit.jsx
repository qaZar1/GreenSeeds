// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";

const BunkerEdit = () => {
    const notify = useNotify();

    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Редактирование бункера"
            mutationOptions={{
                onSuccess: () => notify("Бункер изменён", { type: "success" }),
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />}>
                <NumberInput source="bunker" label="Бункер" disabled/>
                <NumberInput source="distance" label="Расстояние" />
            </SimpleForm>
        </Edit>
    );
};

export default BunkerEdit;
