// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";
import { useNotify } from "react-admin";
import { useRedirect } from "react-admin";

const PlacementCreate = () => {
    const notify = useNotify();
    const redirect = useRedirect();

    return (
        <Create
            sx={{ padding: 2 }}
            actions={<BackButton />}
            mutationMode="pessimistic"
            title="Создание расположения семян"
            mutationOptions={{
                onSuccess: (response) => {
                    notify("Бункер для семян успешно выбран", { type: "success" })
                    redirect('edit', 'placements', response.id);
                },
                onError: (error) => notify("Невозможно указать расположение", { type: "error" }),
            }}
        >
            <SimpleForm toolbar={<ToolbarSave />}>
                <ReferenceInput source="bunker" reference="bunkers" filter={{ free: true }}>
                    <AutocompleteInput optionText="bunker" id="bunker" label="Бункер" noOptionsText="Свободных бункеров нет"/>
                </ReferenceInput>

                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed_ru" id="seed" label="Семена" noOptionsText="Семян нет"/>
                </ReferenceInput>
                <NumberInput source="amount" label="Количество семян в бункере(шт лотков)" min={0} defaultValue={0}/>
            </SimpleForm>
        </Create>
    );
}

export default PlacementCreate;
