// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, AutocompleteInput } from "react-admin";

const AssignmentsCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic">
            <SimpleForm toolbar={<ToolbarSave />}>
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
                <NumberInput source="amount" label="Количество" min={1}/>
            </SimpleForm>
        </Create>
    );
}

export default AssignmentsCreate;
