// SeedCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const validateSeed = (values) => {
    const errors = {};
    if (values.min_density > values.max_density) {
        errors.min_density = "Минимальная плотность не может быть больше максимальной";
        errors.max_density = "Максимальная плотность не может быть меньше минимальной";
    }
    return errors;
};

const SeedCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание семян"
            mutationOptions={{
                onSuccess: (response) => {
                    notify("Семена добавлены", { type: "success" })
                    redirect('edit', 'seeds', response.id);
                },
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />} validate={validateSeed}>
                <TextInput source="seed" label="Семена" />
                <TextInput source="seed_ru" label="Семена (рус)" />
                <NumberInput source="min_density" label="Минимальная плотность (штук/см²)" />
                <NumberInput source="max_density" label="Максимальная плотность (штук/см²)" />
                <NumberInput source="tank_capacity" label="Емкость бункера (кол-во лотков)" />
            </SimpleForm>
        </Create>
    );
}

export default SeedCreate;
