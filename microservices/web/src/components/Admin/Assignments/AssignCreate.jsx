import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { useDataProvider, ReferenceInput, AutocompleteInput } from "react-admin";
import { useEffect, useState } from "react";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const AssignmentsCreate = () => {
    const dataProvider = useDataProvider();
    const [shifts, setShifts] = useState([]);
    const notify = useNotify();
    const redirect = useRedirect();

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
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание сменного задания"
            mutationOptions={{
                onSuccess: (response) => {
                    notify("Сменное задание создано", { type: 'success' })
                    redirect('edit', 'assignments', response.id);
                },
                onError: (error) => {
                    notify(error.message, { type: 'error' });
                }
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <ReferenceInput source="shift" reference="shifts">
                    <AutocompleteInput
                        source="shift"
                        choices={shifts}
                        optionText={(record) => 
                            new Date(record.dt).toLocaleDateString('ru-RU', {
                                year: 'numeric',
                                month: '2-digit',
                                day: '2-digit',
                                hour: 'numeric',
                                minute: 'numeric'
                            })
                        }
                        optionValue="shift"
                        label="Смена"
                    />
                </ReferenceInput>
                <NumberInput source="number" label="Номер сменного задания" min={1}/>
                <ReferenceInput source="receipt" reference="receipts">
                    <AutocompleteInput
                        optionText="description"
                        optionValue="receipt"
                        id="receipt"
                        label="Рецепт"
                    />
                </ReferenceInput>
                <NumberInput source="amount" label="Количество" min={1}/>
            </SimpleForm>
        </Create>
    );
}

export default AssignmentsCreate;
