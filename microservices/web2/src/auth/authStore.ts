import type { AuthData } from "../types/auth";
import { getStoredAuth } from "./authStorage";

let authCache: AuthData | null = getStoredAuth();

export const getAuth = (): AuthData | null => {
  return authCache;
};

export const setAuth = (auth: AuthData | null): void => {
  authCache = auth;
};