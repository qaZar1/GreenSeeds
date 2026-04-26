import { createContext, useContext, useState, type ReactNode } from "react";
import { jwtDecode } from "jwt-decode";
import type { AuthContextType, AuthData, LoginCredentials } from "../types/auth";
import { getStoredAuth, setStoredAuth, clearStoredAuth } from "../auth/authStorage";
import { api } from "../api/apiProvider";

interface Props {
  children: ReactNode;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: Props) => {
  const [auth, setAuth] = useState<AuthData | null>(getStoredAuth());

  const login = async ({ username, password }: LoginCredentials) => {
    const response = await api.create("auth", { username, password });

    console.log(response.data)
    const { access_token } = await response.data;

    const decoded: any = jwtDecode(access_token);

		const decoded_username: string = decoded.username;
    const decoded_user_id: number = decoded.user_id;
    const decoded_full_name: string = decoded.full_name;

    const authData: AuthData = {
      token: access_token,
      role: decoded.role,
      user_id: decoded_user_id,
      username: decoded_username,
    	full_name: decoded_full_name,
    };

    setStoredAuth(authData);
    setAuth(authData);

    return authData;
  };

  const logout = () => {
    clearStoredAuth();
    setAuth(null);
  };

  const value: AuthContextType = {
		auth,
		token: auth?.token ?? null,
		role: auth?.role ?? null,
    user_id: auth?.user_id ?? null,
		username: auth?.username ?? null,
		full_name: auth?.full_name ?? null,
		isAuthenticated: !!auth,
		login,
		logout,
	};

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
};