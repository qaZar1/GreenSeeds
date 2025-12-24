import { getToken } from "../../dataProvider";
import { HttpError } from "react-admin";
import { actionText, handleApiError } from "./api";

export async function apiRequest(url, options = {}) {
    const token = getToken();

    const response = await fetch(url, {
        ...options,
        headers: {
            "Content-Type": "application/json",
            Authorization: token ? `Bearer ${token}` : undefined,
            ...options.headers,
        },
    });

    // 204 → просто успех без тела
    if (response.status === 204) return null;

    // Ошибка → пробуем достать текст/JSON с сервера
    if (!response.ok) {
        throw new HttpError(`HTTP ${response.status}`, response.status, response);
    }

    // Если есть JSON — возвращаем
    try {
        return await response.json();
    } catch (_) {
        return null;
    }
}
