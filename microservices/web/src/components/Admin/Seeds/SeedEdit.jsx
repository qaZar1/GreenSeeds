// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";

const validateSeed = (values) => {
    const errors = {};
    if (values.min_density > values.max_density) {
        errors.min_density = "Минимальная плотность не может быть больше максимальной";
        errors.max_density = "Максимальная плотность не может быть меньше минимальной";
    }
    return errors;
};

const SeedEdit = () => {
    const notify = useNotify();
    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title={"Редактирование семян"}
            mutationOptions={{
                onSuccess: () => notify("Семена изменены", { type: "success" }),
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />} validate={validateSeed}>
                <TextInput source="seed" label="Семена" disabled/>
                <TextInput source="seed_ru" label="Семена (рус)"/>
                <NumberInput source="min_density" label="Минимальная плотность (штук/см²)" />
                <NumberInput source="max_density" label="Максимальная плотность (штук/см²)" />
                <NumberInput source="tank_capacity" label="Емкость бункера (кол-во лотков)" />
            </SimpleForm>
        </Edit>
    );
};

export default SeedEdit;
