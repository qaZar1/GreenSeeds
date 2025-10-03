export const authProvider = {
    // вызывается при логине
    login: async ({ username, password }) => {
      const request = new Request("/api/login", {
        method: "POST",
        body: JSON.stringify({ username, password }),
        headers: new Headers({ "Content-Type": "application/json" }),
      });
  
      const response = await fetch(request);
      if (response.status < 200 || response.status >= 300) {
        throw new Error("Ошибка входа");
      }
  
      const { token, role } = await response.json();
  
      localStorage.setItem("auth", JSON.stringify({ token, role }));
      return Promise.resolve();
    },
  
    // выход
    logout: () => {
      localStorage.removeItem("auth");
      return Promise.resolve();
    },
  
    // проверка авторизации
    checkAuth: () =>
      localStorage.getItem("auth") ? Promise.resolve() : Promise.reject(),
  
    // проверка прав доступа (по ролям)
    getPermissions: () => {
      const auth = localStorage.getItem("auth");
      if (auth) {
        const { role } = JSON.parse(auth);
        return Promise.resolve(role);
      }
      return Promise.reject();
    },
  
    // проверка статуса при ошибках API
    checkError: (error) => {
      const status = error.status;
      if (status === 401 || status === 403) {
        localStorage.removeItem("auth");
        return Promise.reject();
      }
      return Promise.resolve();
    },
  };
  