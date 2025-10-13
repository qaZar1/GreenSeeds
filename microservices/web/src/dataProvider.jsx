// src/dataProvider.js
import { jwtDecode } from "jwt-decode";
import { useNotify } from "react-admin";

export const getToken = () => {
    try {
        const stored = localStorage.getItem("auth");
        if (stored) {
            const parsed = JSON.parse(stored);
            return parsed?.token || null;
        }
    } catch (e) {
        console.warn("Ошибка получения профиля:", e);
    }
    return null;
};

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
        } else if (resource === 'users') {
            switch (action) {
                case 'list': return `/api/users/get`;
                case 'one': return `/api/users/get${idPart}`;
                case 'create': return `/api/register`;
                case 'update': return `/api/users/update`;
                case 'delete': return `/api/users/delete${idPart}`;
            }
        } else if (resource === 'profile'){
            switch (action) {
                case 'one': return `/api/users/get${idPart}`;
                case 'update': return `/api/users/update`;
            }
        } else if (resource === 'placements'){
            switch (action) {
                case 'list': return `/api/placement/get`;
                case 'one': return `/api/placement/get${idPart}`;
                case 'create': return `/api/placement/add`;
                case 'update': return `/api/placement/update`;
                case 'delete': return `/api/placement/delete${idPart}`;
            }
        } else if (resource === 'receipts') {
            switch (action) {
                case 'list': return `/api/receipts/get`;
                case 'one': return `/api/receipts/get${idPart}`;
                case 'create': return `/api/receipts/add`;
                case 'update': return `/api/receipts/update`;
                case 'delete': return `/api/receipts/delete${idPart}`;
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
        } else if (resource === 'users' || resource === 'profile') {
            return {
                ...item,
                id: item.id ?? item.username,
            };
        } else if (resource === 'placements') {
            return {
                ...item,
                id: item.id ?? item.bunker,
            };
        } else if (resource === 'receipts') {
            return {
                ...item,
                id: item.id ?? item.receipt,
            };
        }
        return item;
    },

    // --- Методы CRUD ---

    getList: async (resource, params) => {
        const token = getToken();
        const url = dataProvider.getApiUrl(resource, 'list');
        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            //notify(`Ошибка загрузки данных для ${resource}`, { type: 'error' });
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

    getMany: async (resource, params) => {
        const token = getToken();
        const url = dataProvider.getApiUrl(resource, 'list'); // Используем общий эндпоинт для получения списка

        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            throw new Error(`Ошибка загрузки данных для ${resource}`);
        }

        const data = await response.json();

        // фильтруем только нужные ID
        const filtered = data
            .map(item => dataProvider.transformData(resource, item))
            .filter(item => params.ids.includes(item.id));

        return { data: filtered };
    },

    getOne: async (resource, params) => {
        const token = getToken();
        const url = dataProvider.getApiUrl(resource, 'one', params.id);
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            throw new Error(`Ошибка загрузки данных для ${resource}`);
        }

        const data = await response.json();

        return {
            data: dataProvider.transformData(resource, data),
        };
    },

    create: async (resource, params) => {
        const token = getToken();
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
        } else if (resource === 'users') {
            bodyData = {
                username: params.data.username,
                password: params.data.password,
                full_name: params.data.full_name,
                is_admin: params.data.is_admin,
            };
        } else if (resource === 'placements') {
            bodyData = {
                bunker: params.data.bunker,
                seed: params.data.seed,
            };
        } else if (resource === 'receipts') {
            bodyData = {
                receipt: params.data.receipt,
                seed: params.data.seed,
                gcode: params.data.gcode,
                description: params.data.description,
            };
        } else {
            throw new Error(`Неподдерживаемый ресурс для создания: ${resource}`);
        }

        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(bodyData),
        });

        if (response.status === 204) {
            return {
                data: dataProvider.transformData(resource, params.data),
            };
        }

        if (!response.ok) {
            throw new Error(`Ошибка создания ${resource}`);
        }

        const data = await response.json();

        return {
            data: dataProvider.transformData(resource, data),
        };
    },

    update: async (resource, params) => {
        const token = getToken();
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
        } else if (resource === 'users') {
            bodyData = {
                username: params.data.username,
                password: params.data.password,
                full_name: params.data.full_name,
                is_admin: params.data.is_admin,
            };
        } else if (resource === 'profile') {
            bodyData = {
                username: params.data.username,
                full_name: params.data.full_name,
            };
        } else if (resource === 'placements') {
            bodyData = {
                bunker: params.data.id,
                seed: params.data.seed,
            };
        } else if (resource === 'receipts') {
            bodyData = {
                receipt: params.data.id,
                seed: params.data.seed,
                gcode: params.data.gcode,
                description: params.data.description,
            };
        } else {
            throw new Error(`Неподдерживаемый ресурс для обновления: ${resource}`);
        }

        const response = await fetch(url, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(bodyData),
        });        

        if (response.status === 204) {
            return {
                data: dataProvider.transformData(resource, params.data),
            };
        }

        if (!response.ok) {
            throw new Error(`Ошибка обновления ${resource}`);
        }

        const data = await response.json();

        return {
            data: dataProvider.transformData(resource, data),
        };
    },

    delete: async (resource, params) => {
        const token = getToken();
        const url = dataProvider.getApiUrl(resource, 'delete', params.id);
        const response = await fetch(url, {
            method: "DELETE",
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            throw new Error(`Ошибка удаления ${resource}`);
        }

        return { data: params.previousData };
    },
};

export default dataProvider;