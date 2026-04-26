import { apiClient } from "./apiClient";
import { resourcesConfig } from "./resourcesConfig";

const getApiUrl = (resource: string, action: string, params?: any) => {
  const cfg = resourcesConfig[resource as keyof typeof resourcesConfig];

  if (!cfg) throw new Error(`Unknown resource ${resource}`);

  const url = cfg.urls[action];

  if (!url) throw new Error(`Unknown action ${action}`);

  const baseUrl =
    typeof url === "function"
      ? url(params)
      : url;

  if (action === "list" && params && typeof params === "object") {
    const query = new URLSearchParams(params).toString();
    return query ? `${baseUrl}?${query}` : baseUrl;
  }

  return baseUrl;
};

const transformData = (resource: string, item: any) => {
  if (!item) return item;

  const key = resourcesConfig[resource]?.idKey;

  if (!key) return item;

  return {
    ...item,
    id: item.id ?? item[key],
  };
};

export const api = {
  async getList(resource: string, params?: any) {
    const url = getApiUrl(resource, "list", params);
    const res = await apiClient(url);
    const data = res?.data;

    if (!data) {
      return [];
    }

    return data.map((item: any) =>
      transformData(resource, item)
    );
  },

  async getOne(resource: string, id: any) {
    const url = getApiUrl(resource, "one", id);

    const data = await apiClient(url);

    return transformData(resource, data);
  },

  async create(resource: string, data: any, options: RequestInit = {}) {
    const cfg = resourcesConfig[resource];

    const url = getApiUrl(resource, "create", data);

    const payload = cfg?.createPayload
      ? cfg.createPayload(data)
      : data;

    const res = await apiClient(url, {
      method: "POST",
      body: JSON.stringify(payload),
      ...options,
    });

    return transformData(resource, res ?? data);
  },

  async update(resource: string, data: any) {
    const cfg = resourcesConfig[resource];

    const url = getApiUrl(resource, "update");

    const payload = cfg.updatePayload
      ? cfg.updatePayload(data)
      : data;

    const res = await apiClient(url, {
      method: "PUT",
      body: JSON.stringify(payload),
    });

    return transformData(resource, res ?? data);
  },

  async delete(resource: string, id: any) {
    const url = getApiUrl(resource, "delete", id);

    await apiClient(url, {
      method: "DELETE",
    });

    return true;
  },
};