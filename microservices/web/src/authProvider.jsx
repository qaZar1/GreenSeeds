import { jwtDecode } from "jwt-decode";

export const authProvider = {
  login: async ({ username, password }) => {
    const response = await fetch("/api/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
      headers: { "Content-Type": "application/json" },
    });

    if (response.status < 200 || response.status >= 300) {
      throw new Error("Ошибка входа");
    }

    const { access_token } = await response.json();

    // Раскодируем JWT
    const decoded = jwtDecode(access_token);

    const role = decoded.role; // в твоем примере "admin"

    localStorage.setItem(
      "auth",
      JSON.stringify({ token: access_token, role })
    );
  },

  logout: () => {
    localStorage.removeItem("auth");
    return Promise.resolve();
  },

  checkAuth: () =>
    localStorage.getItem("auth") ? Promise.resolve() : Promise.reject(),

  getPermissions: () => {
    const auth = localStorage.getItem("auth");
    if (!auth) return Promise.reject();
    const { role } = JSON.parse(auth);
    return Promise.resolve(role);
  },

  checkError: (error) => {
    if (error.status === 401 || error.status === 403) {
      localStorage.removeItem("auth");
      return Promise.reject();
    }
    return Promise.resolve();
  },
};
