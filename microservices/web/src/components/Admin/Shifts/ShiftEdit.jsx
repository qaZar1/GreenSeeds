import React from "react";
import { Edit, SimpleForm, DateTimeInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useNotify } from "react-admin";

const ShiftEdit = () => {
    const notify = useNotify();

    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Редактирование смены"
            mutationOptions={{
                onSuccess: () => notify("Смена изменена", { type: "success" }),
                onError: (error) => notify(error.message, { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />}>
                <DateTimeInput source="dt" label="Дата" />
            </SimpleForm>
        </Edit>
    );
};

export default ShiftEdit;