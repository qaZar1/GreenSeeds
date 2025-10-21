// SeedCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";

const SeedCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarSave />}>
                <TextInput source="seed" label="Семена" />
                <NumberInput source="min_density" label="Минимальная плотность" />
                <NumberInput source="max_density" label="Максимальная плотность" />
                <NumberInput source="tank_capacity" label="Емкость бункера" />
                <NumberInput source="latency" label="Задержка" />
            </SimpleForm>
        </Create>
    );
}

export default SeedCreate;
