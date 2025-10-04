// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput, useNotify } from "react-admin";
import { ToolbarSave } from "./Toolbars";
import BackButton from "../utils/Back";

const BunkerCreate = () => {
    const notify = useNotify();
    
    const handleSubmit = (data) => {
        // Проверка пустых полей
        if (!data.bunker || !data.distance) {
            notify("Все поля должны быть заполнены", { type: "warning" });
            return;
        }

        // Вывод изменённых данных
        console.log("Изменённые значения:", data);
        alert(`Изменённые значения:\nБункер: ${data.bunker}\nРасстояние: ${data.distance}`);
    };

    return (
    <>
    <Create sx={{ padding: 2 }} actions={<BackButton />}>
        <SimpleForm onSubmit={handleSubmit} toolbar={<ToolbarSave />}>
            <TextInput source="bunker" label="Бункер" />
            <NumberInput source="distance" label="Расстояние" />
        </SimpleForm>
    </Create>
    </>
    );
}

export default BunkerCreate;
