import type { AuthData } from "../types/auth";

const AUTH_KEY = "auth";

export const getStoredAuth = (): AuthData | null => {
  try {
    const data = localStorage.getItem(AUTH_KEY);
    return data ? JSON.parse(data) : null;
  } catch {
    return null;
  }
};

export const setStoredAuth = (auth: AuthData): void => {
  localStorage.setItem(AUTH_KEY, JSON.stringify(auth));
};

export const clearStoredAuth = (): void => {
  localStorage.removeItem(AUTH_KEY);
};