import React from "react";
import { useDataProvider } from "react-admin";
import { useEffect, useState } from "react";
import { AutocompleteInput } from "react-admin";

const ShiftAutocomplete = ({ source, ...props }) => { // <- добавили source
    const dataProvider = useDataProvider();
    const [choices, setChoices] = useState([]);

    useEffect(() => {
        dataProvider.getList('shifts', { pagination: { page: 1, perPage: 1000 } })
            .then(({ data }) => {
                const today = new Date();
                const filtered = data.filter(shift => new Date(shift.dt) >= today);
                const sorted = filtered.sort((a, b) => new Date(a.dt) - new Date(b.dt));
                setChoices(sorted);
            });
    }, [dataProvider]);

    return <AutocompleteInput source={source} {...props} choices={choices} />; // <- передаем source
};

export default ShiftAutocomplete;
