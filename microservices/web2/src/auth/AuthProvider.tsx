import { createContext, useContext, useState, type ReactNode, useMemo } from "react";
import { jwtDecode } from "jwt-decode";
import type { AuthContextType, AuthData, LoginCredentials } from "../types/auth";

import {
  getStoredAuth,
  setStoredAuth,
  clearStoredAuth,
} from "./authStorage";

import { getAuth, setAuth as setAuthStore } from "./authStore";

interface Props {
  children: ReactNode;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: Props) => {
  const [auth, setAuth] = useState<AuthData | null>(getAuth());

  const login = async ({ username, password }: LoginCredentials) => {
    const response = await fetch("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
      headers: { "Content-Type": "application/json" },
    });

    if (!response.ok) {
      throw new Error("Login failed");
    }

    const { access_token } = await response.json();

    const decoded: any = jwtDecode(access_token);
    const decoded_username: string = decoded.username;

    const responseFullName = await fetch("/api/users/get/" + decoded_username, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + access_token,
      },
    });

    if (!responseFullName.ok) {
      throw new Error("Login failed");
    }

    const data = await responseFullName.json();

    const authData: AuthData = {
      user_id: decoded.user_id,
      token: access_token,
      role: decoded.role,
      username: decoded_username,
      full_name: data.full_name,
    };

    // persist
    setStoredAuth(authData);

    // memory cache
    setAuthStore(authData);

    // react state
    setAuth(authData);

    return authData;
  };

  const logout = () => {
    clearStoredAuth();
    setAuthStore(null);
    setAuth(null);
  };

  const value: AuthContextType = useMemo(
    () => ({
      auth,
      token: auth?.token ?? null,
      role: auth?.role ?? null,
      user_id: auth?.user_id ?? null,
      username: auth?.username ?? null,
      full_name: auth?.full_name ?? null,
      isAuthenticated: !!auth,
      login,
      logout,
    }),
    [auth]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
};