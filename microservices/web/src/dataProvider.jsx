import { jwtDecode } from "jwt-decode";
import { useNotify } from "react-admin";
import { fetchUtils } from 'react-admin';

const LIMIT = 50;

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

const baseProvider = {
    /* Возвращает полный URL для API-запроса */
    getApiUrl: (resource, action, id) => {
        const idPart = id !== undefined && id !== null ? `/${id}` : '';

        if (resource === 'bunkers') {
            switch (action) {
                case 'list': return `/api/bunkers/get`;
                case 'one': return `/api/bunkers/get${idPart}`;
                case 'create': return `/api/bunkers/add`;
                case 'update': return `/api/bunkers/update`;
                case 'delete': return `/api/bunkers/delete${idPart}`;
            }
        } else if (resource === 'seeds') {
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
        } else if (resource === 'shifts') {
            switch (action) {
                case 'list': return `/api/shifts/get`;
                case 'one': return `/api/shifts/get${idPart}`;
                case 'create': return `/api/shifts/add`;
                case 'update': return `/api/shifts/update`;
                case 'delete': return `/api/shifts/delete${idPart}`;
            }
        } else if (resource === 'assignments') {
            switch (action) {
                case 'list': return `/api/assignments/get`;
                case 'one': return `/api/assignments/get${idPart}`;
                case 'create': return `/api/assignments/add`;
                case 'update': return `/api/assignments/update`;
                case 'delete': return `/api/assignments/delete${idPart}`;
            }
        } else if (resource === 'reports') {
            switch (action) {
                case 'list': return `/api/reports/get`;
                case 'one': return `/api/reports/get${idPart}`;
            }
        } else if (resource === 'choice') {
            switch (action) {
                case 'list': return `/api/shifts/getWithoutUser`;
                case 'update': return `/api/shifts/update`;
            }
        } else if (resource === 'tasks') {
            switch (action) {
                case 'list': return `/api/assignments/active-tasks${idPart}`;
            }
        } else if (resource === 'task') {
            switch (action) {
                case 'one': return `/api/assignments/task${idPart}`;
            }
        } else if (resource === 'logs') {
            switch (action) {
                case 'list': return `/api/logs/get`;
            }
        }
        throw new Error(`Неподдерживаемый ресурс или действие: ${resource}/${action}`);
    },

    /* Добавляет поле 'id' к элементу */
    transformData: (resource, item) => {
        if (!item) return item;

        if (resource === 'bunkers') {
            return {
                ...item,
                id: item.id ?? item.bunker,
            };
        } else if (resource === 'seeds') {
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
        } else if (resource === 'shifts') {
            return {
                ...item,
                id: item.id ?? item.shift,
            };
        } else if (resource === 'assignments') {
            return {
                ...item,
                id: item.id ?? item.id,
            };
        } else if (resource === 'reports') {
            return {
                ...item,
                id: item.id ?? item.id,
            };
        } else if (resource === 'choice') {
            return {
                ...item,
                id: item.id ?? item.shift,
            };
        } else if (resource === 'tasks') {
            return {
                ...item,
                id: item.id ?? item.username,
            };
        } else if (resource === 'task') {
            return {
                ...item,
                id: item.id ?? item.id,
            };
        } else if (resource === 'logs') {
            return {
                ...item,
                id: item.id ?? item.id,
            }
        }
        return item;
    },

    /* Методы CRUD */
    getList: async (resource, params) => {
        const token = getToken();
        
        let url
        if (resource === 'tasks'){
            url = dataProvider.getApiUrl(resource, 'list', params.filter.username);
        } else if (resource === 'logs'){
            const search = params.filter?.search || "";
            const level = params.filter?.level || "ALL";
            const page = params.pagination?.page || 1;
            const perPage = params.pagination?.perPage || 50;
            const offset = (page - 1) * perPage;
            const dateFrom = params.filter?.date_from || "";
            const dateTo = params.filter?.date_to || ""

            const qs = new URLSearchParams({
                search,
                level,
                limit: perPage,
                offset,
                date_from: dateFrom,
                date_to: dateTo,
            });

            url = `/api/logs/get?${qs.toString()}`;
        } else {
            url = dataProvider.getApiUrl(resource, 'list');
        }

        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        });
    
        if (response.status === 404){
            return {
                data: [],
                total: 0,
            };
        } else if (!response.ok) {
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
        const url = dataProvider.getApiUrl(resource, 'list');

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
            bodyData = {
                seed: params.data.seed,
                seed_ru: params.data.seed_ru,
                min_density: params.data.min_density,
                max_density: params.data.max_density,
                tank_capacity: params.data.tank_capacity,
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
                seed: params.data.seed,
                gcode: params.data.gcode,
                description: params.data.description,
            };
        } else if (resource === 'shifts') {
            const dt = params.data.dt;
            const date = new Date(dt);
            bodyData = {
                dt: date.toISOString(),
            };
        } else if (resource === 'assignments') {
            bodyData = {
                shift: params.data.shift,
                number: params.data.number,
                receipt: params.data.receipt,
                amount: params.data.amount,
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
                bunker: params.data.id,
                distance: params.data.distance,
            };
        } else if (resource === 'seeds') {
            bodyData = {
                seed: params.data.id,
                seed_ru: params.data.seed_ru,
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
                seed_ru: params.data.seed_ru,
            };
        } else if (resource === 'receipts') {
            bodyData = {
                receipt: params.data.id,
                seed: params.data.seed,
                gcode: params.data.gcode,
                description: params.data.description,
            };
        } else if (resource === 'shifts') {
            const dt = params.data.dt;
            const date = new Date(dt);
            bodyData = {
                shift: params.data.id,
                dt: date.toISOString(),
            };
        } else if (resource === 'assignments') {
            bodyData = {
                id: params.data.id,
                shift: params.data.shift,
                number: params.data.number,
                receipt: params.data.receipt,
                amount: params.data.amount,
            };
        } else if (resource === 'tasks') {
            bodyData = {
                shift: params.data.id,
                username: params.data.username,
                dt: params.data.dt,
            };
        } else if (resource === 'choice') {
            bodyData = {
                shift: params.data.id,
                username: params.data.username,
                dt: params.data.dt,
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
        const token = getToken()
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

const dataProvider = {
    ...baseProvider,
    getList: async (resource, params) => {
      if (resource === "bunkers" && params.filter?.free) {
        try {
            if (resource === "bunkers" && params.filter?.free) {
              const response = await fetch("/api/bunkers/getForPlacement", {
                headers: {
                  'Authorization': `Bearer ${getToken()}`,
                },
              });
              if (!response.ok) {
                return { data: [], total: 0 };
              }

              const json = await response.json();
              
              return { data: json.map((item) => ({
                ...item,
                id: item.id ?? item.bunker,
              })), total: json.length };
            }
            return baseProvider.getList(resource, params);
          } catch (e) {
            console.error("Error in getList:", e);
            throw e;
          }
      }
      return baseProvider.getList(resource, params);
    },
  };
  

export default dataProvider;