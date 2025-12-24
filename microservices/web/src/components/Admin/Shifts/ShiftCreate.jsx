import React from "react";
import { Create, SimpleForm, DateTimeInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const ShiftCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание смены"
            mutationOptions={{
                onSuccess: (response) => {
                    notify("Смена создана", { type: "success" });
                    redirect('edit', 'shifts', response.id);
                },
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <DateTimeInput source="dt" label="Дата"/>
            </SimpleForm>
        </Create>
    );
}

export default ShiftCreate;
