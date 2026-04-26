export interface AuthData {
  token: string;
  role: string;
  user_id: number;
  username: string;
  full_name: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface AuthContextType {
  auth: AuthData | null;
  token: string | null;
  role: string | null;
  user_id: number | null;
  username: string | null;
  full_name: string | null;
  isAuthenticated: boolean;
  login: (credentials: LoginCredentials) => Promise<AuthData>;
  logout: () => void;
}