// BunkerEdit.jsx
import React from "react";
import { Edit, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarEdit } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";

const AssignmentsEdit = () => {
    return (
        <Edit sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic" title="Редактирование сменного задания">
            <SimpleForm toolbar={<ToolbarEdit />}>
                <ReferenceInput source="shift" reference="shifts">
                    <AutocompleteInput
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
                        id="shift"
                        label="Смена"
                    />
                </ReferenceInput>
                <NumberInput source="number" label="Номер сменного задания"/>
                <ReferenceInput source="receipt" reference="receipts">
                    <AutocompleteInput
                        optionText="description"
                        optionValue="receipt"
                        id="receipt"
                        label="Рецепт"
                    />
                </ReferenceInput>
                <NumberInput source="amount" label="Количество"/>
            </SimpleForm>
        </Edit>
    );
};

export default AssignmentsEdit;
