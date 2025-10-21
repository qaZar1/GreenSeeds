// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";

const SeedEdit = () => {
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <TextInput source="seed" label="Семена" disabled/>
                <NumberInput source="min_density" label="Минимальная плотность" />
                <NumberInput source="max_density" label="Максимальная плотность" />
                <NumberInput source="tank_capacity" label="Емкость бункера" />
                <NumberInput source="latency" label="Задержка" />
            </SimpleForm>
        </Edit>
    );
};

export default SeedEdit;
