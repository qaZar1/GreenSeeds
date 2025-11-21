// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";

const BunkerCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic" title="Создание бункера">
            <SimpleForm toolbar={<ToolbarSave />}>
                <NumberInput source="bunker" label="Бункер" max={10} min={1}/>
                <NumberInput source="distance" label="Расстояние" />
            </SimpleForm>
        </Create>
    );
}

export default BunkerCreate;
