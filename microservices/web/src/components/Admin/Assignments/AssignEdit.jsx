import React from "react";
import { useDataProvider } from "react-admin";
import { useEffect, useState } from "react";
import { AutocompleteInput } from "react-admin";
import { Edit, SimpleForm, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput } from "react-admin";
import { useNotify } from "react-admin";

const AssignmentsEdit = () => {
    const dataProvider = useDataProvider();
    const [shifts, setShifts] = useState([]);
    const notify = useNotify();

    useEffect(() => {
        dataProvider
            .getList("shifts")
            .then(({ data }) => {
                const today = new Date();
                today.setHours(0, 0, 0, 0);
                const filtered = data.filter((shift) => new Date(shift.dt) >= today);
                const sorted = filtered.sort((a, b) => new Date(a.dt) - new Date(b.dt));
                setShifts(sorted);
            });
    }, [dataProvider]);

    return (
        <Edit
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Редактирование сменного задания"
            mutationOptions={{
                onSuccess: () => {
                    notify("Сменное задание изменено", { type: 'success' })
                },
                onError: (error) => {
                    notify(error.message, { type: 'error' });
                }
            }}
        >
            <SimpleForm toolbar={<ToolbarEdit />}>
                <AutocompleteInput
                    source="shift"
                    choices={shifts}
                    optionText={(record) =>
                        new Date(record.dt).toLocaleDateString("ru-RU", {
                            year: "numeric",
                            month: "2-digit",
                            day: "2-digit",
                            hour: "numeric",
                            minute: "numeric",
                        })
                    }
                    optionValue="shift"
                    id="shift"
                    label="Смена"
                />
                <NumberInput source="number" label="Номер сменного задания" />
                <ReferenceInput source="receipt" reference="receipts">
                    <AutocompleteInput optionText="description" optionValue="receipt" label="Рецепт" />
                </ReferenceInput>
                <NumberInput source="amount" label="Количество" />
            </SimpleForm>
        </Edit>
    );
};

export default AssignmentsEdit;
