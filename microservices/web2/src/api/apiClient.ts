import toast from "react-hot-toast";
import { getStoredAuth } from "../auth/authStorage";

export const apiClient = async (
  url: string,
  options: RequestInit = {}
) => {
  const auth = getStoredAuth();

  const headers = new Headers(options.headers);

  headers.set("Content-Type", "application/json");

  if (auth?.token) {
    headers.set("Authorization", `Bearer ${auth.token}`);
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem("auth");
      window.location.href = "/login";
      return;
    }

    if (response.status === 403) {
      toast.error("Нет прав доступа");
      window.location.replace("/");
      throw new Error("FORBIDDEN");
    }

    if (response.status === 404) {
      return {
        data: [],
        headers: response.headers,
      };
    }

    if (response.status >= 500) {
      throw new Error("SERVER_ERROR");
    }

    throw new Error(`API error ${response.status}`);
  }

  let data = null;

  if (response.status !== 204) {
    const contentType = response.headers.get("content-type") || "";

    if (contentType.includes("application/json")) {
      data = await response.json();
    } else {
      data = await response.blob();
    }
  }
  
  return {
    data,
    headers: response.headers,
  };
};