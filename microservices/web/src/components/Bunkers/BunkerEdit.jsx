// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput, useNotify } from "react-admin";
import { ToolbarEdit } from "./Toolbars";
import BackButton from "../utils/Back";
import { Box } from "@mui/material";

const BunkerEdit = () => {
    const notify = useNotify();

    const handleSubmit = (data) => {
        // Проверка пустых полей
        if (!data.distance) {
            notify("Все поля должны быть заполнены", { type: "warning" });
            return;
        }

        // Вывод изменённых данных
        console.log("Изменённые значения:", data);
        alert(`Изменённые значения:\nРасстояние: ${data.distance}`);
    };

    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />}>
            <SimpleForm onSubmit={handleSubmit} toolbar={<ToolbarEdit />}>
                <TextInput source="bunker" label="Бункер" disabled/>
                <NumberInput source="distance" label="Расстояние" />
            </SimpleForm>
        </Edit>
    );
};

export default BunkerEdit;
