// BunkerCreate.jsx
import React from "react";
import { Create, SimpleForm, TextInput, NumberInput } from "react-admin";
import { ToolbarSave } from "../../utils/Toolbars";
import BackButton from "../../utils/Back";
import { ReferenceInput, SelectInput, AutocompleteInput } from "react-admin";

const PlacementCreate = () => {
    return (
        <Create sx={{ padding: 2 }} actions={<BackButton />} mutationMode="pessimistic" title="Создание расположения семян">
            <SimpleForm toolbar={<ToolbarSave />}>
                <ReferenceInput source="bunker" reference="bunkers" filter={{ free: true }}>
                    <AutocompleteInput optionText="bunker" id="bunker" label="Бункер" noOptionsText="Свободных бункеров нет"/>
                </ReferenceInput>

                <ReferenceInput source="seed" reference="seeds">
                    <AutocompleteInput optionText="seed_ru" id="seed" label="Семена" noOptionsText="Семян нет"/>
                </ReferenceInput>
            </SimpleForm>
        </Create>
    );
}

export default PlacementCreate;
