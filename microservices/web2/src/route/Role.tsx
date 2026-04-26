import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export const RequireRole = ({ role }: { role: "admin" | "operator" }) => {
  const { role: userRole } = useAuth();

  if (userRole !== role) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
};