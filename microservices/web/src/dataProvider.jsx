// src/dataProvider.js
const dataProvider = {
    // --- Вспомогательные функции для определения API и трансформации данных ---

    /**
     * Возвращает полный URL для API-запроса
     */
    getApiUrl: (resource, action, id = '') => {
        const idPart = id ? `/${id}` : '';

        if (resource === 'bunkers') {
            switch (action) {
                case 'list': return `/api/bunkers/get`;
                case 'one': return `/api/bunkers/get${idPart}`;
                case 'create': return `/api/bunkers/add`;
                case 'update': return `/api/bunkers/update`;
                case 'delete': return `/api/bunkers/delete${idPart}`;
            }
        } else if (resource === 'seeds') {
            // Предполагаемые эндпоинты для ресурса 'seeds'
            switch (action) {
                case 'list': return `/api/seeds/get`;
                case 'one': return `/api/seeds/get${idPart}`;
                case 'create': return `/api/seeds/add`;
                case 'update': return `/api/seeds/update`;
                case 'delete': return `/api/seeds/delete${idPart}`;
            }
        }
        throw new Error(`Неподдерживаемый ресурс или действие: ${resource}/${action}`);
    },

    /**
     * Добавляет поле 'id' к элементу, используя уникальный ключ ресурса, если 'id' отсутствует
     */
    transformData: (resource, item) => {
        if (!item) return item;

        if (resource === 'bunkers') {
            // Для 'bunkers' уникальный ключ — 'bunker'
            return {
                ...item,
                id: item.id ?? item.bunker,
            };
        } else if (resource === 'seeds') {
            // Для 'seeds' уникальный ключ — 'seed'
            return {
                ...item,
                id: item.id ?? item.seed,
            };
        }
        return item;
    },

    // --- Методы CRUD ---

    getList: async (resource, params) => {
        const url = dataProvider.getApiUrl(resource, 'list');
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`Ошибка загрузки данных для ${resource}`);
        }

        const data = await response.json();

        const dataWithId = data.map(item =>
            dataProvider.transformData(resource, item)
        );

        return {
            data: dataWithId,
            total: dataWithId.length,
        };
    },

    getOne: async (resource, params) => {
        const url = dataProvider.getApiUrl(resource, 'one', params.id);
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`Ошибка загрузки данных для ${resource}`);
        }

        const data = await response.json();

        return {
            data: dataProvider.transformData(resource, data),
        };
    },

    create: async (resource, params) => {
        const url = dataProvider.getApiUrl(resource, 'create');

        let bodyData = {};
        if (resource === 'bunkers') {
            bodyData = {
                bunker: params.data.bunker,
                distance: params.data.distance,
            };
        } else if (resource === 'seeds') {
            // ** Использование требуемой структуры JSON для 'seeds' **
            bodyData = {
                seed: params.data.seed,
                min_density: params.data.min_density,
                max_density: params.data.max_density,
                tank_capacity: params.data.tank_capacity,
                latency: params.data.latency,
            };
        } else {
            throw new Error(`Неподдерживаемый ресурс для создания: ${resource}`);
        }

        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(bodyData),
        });

        if (!response.ok) {
            throw new Error(`Ошибка создания ${resource}`);
        }

        const data = await response.json();

        return {
            data: dataProvider.transformData(resource, data),
        };
    },

    update: async (resource, params) => {
        const url = dataProvider.getApiUrl(resource, 'update');

        let bodyData = {};
        if (resource === 'bunkers') {
            bodyData = {
                bunker: params.data.id, // ID используется как 'bunker' для обновления
                distance: params.data.distance,
            };
        } else if (resource === 'seeds') {
            // ** Использование требуемой структуры JSON для 'seeds' **
            bodyData = {
                seed: params.data.id, // ID используется как 'seed' для обновления
                min_density: params.data.min_density,
                max_density: params.data.max_density,
                tank_capacity: params.data.tank_capacity,
                latency: params.data.latency,
            };
        } else {
            throw new Error(`Неподдерживаемый ресурс для обновления: ${resource}`);
        }

        const response = await fetch(url, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(bodyData),
        });

        if (!response.ok) {
            throw new Error(`Ошибка обновления ${resource}`);
        }

        const data = await response.json();

        return {
            data: dataProvider.transformData(resource, data),
        };
    },

    delete: async (resource, params) => {
        const url = dataProvider.getApiUrl(resource, 'delete', params.id);
        const response = await fetch(url, {
            method: "DELETE",
        });

        if (!response.ok) {
            throw new Error(`Ошибка удаления ${resource}`);
        }

        return { data: params.previousData };
    },
};

export default dataProvider;