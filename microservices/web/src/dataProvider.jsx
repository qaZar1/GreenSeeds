import { jwtDecode } from "jwt-decode";
import { Notification, useNotify } from "react-admin";
import { fetchUtils } from 'react-admin';
import { handleApiError } from "./components/utils/api";
import { resourcesConfig } from "./components/utils/route";
import { apiRequest } from "./components/utils/request";

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

const getApiUrl = (resource, action, id) => {
    const cfg = resourcesConfig[resource];
    if (!cfg) throw new Error(`Unknown resource: ${resource}`);

    const url = cfg.urls[action];
    if (!url) throw new Error(`Unknown action: ${action} for resource ${resource}`);

    return typeof url === "function" ? url(id) : url;
};

const transformData = (resource, item) => {
    if (!item) return item;

    const key = resourcesConfig[resource]?.idKey;
    if (!key) return item;

    return {
        ...item,
        id: item.id ?? item[key],
    };
};

const baseProvider = {
    getList: async (resource, params) => {
        let url
        if (resource === 'tasks'){
            url = getApiUrl(resource, 'list', params.filter.username);
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
        } else if (resource === 'seedWithBunkers') {
            url = getApiUrl(resource, 'list', params.id);
        } else {
            url = getApiUrl(resource, 'list');
        }

        const data = await apiRequest(url);
    
        if (!data){
            return {
                data: [],
                total: 0,
            };
        }
    
    
        const dataWithId = data.map(item =>
            transformData(resource, item)
        );
    
        return {
            data: dataWithId,
            total: dataWithId.length,
        };
    },
    

    getMany: async (resource, params) => {
        const url = getApiUrl(resource, 'list');
        const data = await apiRequest(url);

        // фильтруем только нужные ID
        const filtered = data
            .map(item => transformData(resource, item))
            .filter(item => params.ids.includes(item.id));

        return { data: filtered };
    },

    getOne: async (resource, params) => {
        const url = getApiUrl(resource, 'one', params.id);
        const data = await apiRequest(url);

        return {
            data: transformData(resource, data),
        };
    },

    create: async (resource, params) => {
        const url = getApiUrl(resource, "create");
        const cfg = resourcesConfig[resource];

        if (!cfg.createPayload) {
            throw new Error(`createPayload is not defined for ${resource}`);
        }

        let bodyData = null;
        let data = null;

        if (params !== undefined) {
            bodyData = cfg.createPayload(params.data);
            data = await apiRequest(url, {
                method: "POST",
                body: JSON.stringify(bodyData),
            });
        } else {
            data = await apiRequest(url, {
                method: "POST",
            });
        }

        const response = data ?? params.data;
        console.log("response", response)

        const transformedData = transformData(resource, response)
        console.log("transformedData", transformedData)

        return { data: transformedData };
    },


    update: async (resource, params) => {
        const url = getApiUrl(resource, "update");
        const cfg = resourcesConfig[resource];

        console.log("params.data", params.data)

        if (!cfg.updatePayload) {
            throw new Error(`updatePayload is not defined for ${resource}`);
        }

        const bodyData = cfg.updatePayload(params.data);

        console.log("bd", bodyData)

        const data = await apiRequest(url, {
            method: "PUT",
            body: JSON.stringify(bodyData),
        });

        const response = data ?? params.data;

        return { data: transformData(resource, response) };
    },


    delete: async (resource, params) => {
        const url = getApiUrl(resource, 'delete', params.id);
        await apiRequest(url, { method: "DELETE" });

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