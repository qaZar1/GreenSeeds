// src/dataProvider.js
const dataProvider = {
    getList: async (resource, params) => {
        const response = await fetch(`/api/bunker/get`);
        if (!response.ok) {
            throw new Error("Ошибка загрузки данных");
        }
    
        const data = await response.json();
    
        const dataWithId = data.map(item => ({
            ...item,
            id: item.id ?? item.bunker,
        }));
    
        return {
            data: dataWithId,
            total: dataWithId.length,
        };
    },    
    getOne: async (resource, params) => {
        const response = await fetch(`/api/bunker/get/${params.id}`);
        if (!response.ok) {
            throw new Error("Ошибка загрузки");
        }

        const data = await response.json();

        return {
            data: {
                ...data,
                id: data.id ?? data.bunker,
            },
        };
    },
    create: async (resource, params) => {
        const response = await fetch(`/api/bunker/add`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(
                {
                    bunker: params.data.bunker,
                    distance: params.data.distance,
                }
            ),
        });

        if (!response.ok) {
            throw new Error("Ошибка создания");
        }

        const data = await response.json();

        return {
            data: {
                ...data,
                id: data.id ?? data.bunker,
            },
        };
    },

    update: async (resource, params) => {
        const response = await fetch(`/api/bunker/update`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(
                {
                    bunker: params.data.id,
                    distance: params.data.distance,
                }
            ),
        });

        if (!response.ok) {
            throw new Error("Ошибка обновления");
        }

        const data = await response.json();

        return {
            data: {
                ...data,
                id: data.id ?? data.bunker,
            },
        };
    },

    delete: async (resource, params) => {
        const response = await fetch(`/api/bunker/delete/${params.id}`, {
            method: "DELETE",
        });

        if (!response.ok) {
            throw new Error("Ошибка удаления");
        }

        return { data: params.previousData };
    },
};

export default dataProvider;
