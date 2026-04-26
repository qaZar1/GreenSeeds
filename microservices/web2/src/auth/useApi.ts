import { getAuth } from "../auth/authStore";
import { clearStoredAuth } from "../auth/authStorage";

export const api = async (url: string, options: RequestInit = {}) => {
  const auth = getAuth();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (auth?.token) {
    headers.Authorization = `Bearer ${auth.token}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (response.status === 401 || response.status === 403) {
    clearStoredAuth();
    window.location.href = "/login";
    throw new Error("Unauthorized");
  }

  return response;
};