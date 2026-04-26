import { Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "../components/pages/LoginPage";
import AppLayout from "../components/layout/DashboardLayout";
import UsersPage from "../components/Admin/Users/UserPage";
import ShiftsPage from "../components/Admin/Shifts/ShiftsPage";
import AssignmentsPage from "../components/Admin/Assignment/Assignment";
import ReportsPage from "../components/Admin/Reports/ReportsPage";
import ReportShow from "../components/Admin/Reports/ReportShow";
import LogsPage from "../components/Admin/Logs/LogsPage";
import { ProtectedRoute } from "./ProtectedRoute";
import ProfilePage from "../components/Menu/Profile/ProfilePage";
import DashboardPage from "../components/Admin/Dashboard/Dashboard";
import BunkersPage from "../components/Admin/Bunkers/Bunkers";
import CalibrationPage from "../components/Admin/Calibration/Calibrate";
import DeviceSettingCreatePage from "../components/Admin/DeviceSettings/DeviceSettingsCreatePage";
import DeviceSettingsPage from "../components/Admin/DeviceSettings/DeviceSettingsPage";
import ReceiptPage from "../components/Admin/Receipt/Receipts";
import PlacementPage from "../components/Admin/Placement/Placement";
import SeedsPage from "../components/Admin/Seeds/SeedsPage";
import ReceiptCreatePage from "../components/Admin/Receipt/ReceiptCreatePage";
import SproutLoader from "../components/utils/Loader/SproutLoader";
import NotFoundPage from "../components/pages/NotFoundPage";
import ChoicePage from "../components/Operator/ChoiceTask/ChoicePage";
import AppliedTasksPage from "../components/Operator/AppliedTasks/AppliedTaskPage";
import { RequireRole } from "./Role";
import TaskDetails from "../components/Operator/RunTask/TaskDetails";

export const AppRouter = () => {
  return (
    <Routes>

      <Route path="/login" element={<LoginPage />} />

      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route element={<RequireRole role="admin" />}>
          <Route path="/users" element={<UsersPage />} />
          <Route path="/shifts" element={<ShiftsPage/>} />
          <Route path="/assignments" element={<AssignmentsPage/>} />
          <Route path="/reports" element={<ReportsPage/>} />
          <Route path="/reports/:id" element={<ReportShow/>} />
          <Route path="/logs" element={<LogsPage/>} />

          <Route path="/settings/bunkers" element={<BunkersPage/>} />
          <Route path="/settings/seeds" element={<SeedsPage/>} />
          <Route path="/settings/placements" element={<PlacementPage/>} />
          <Route path="/settings/receipts" element={<ReceiptPage/>} />
          <Route path="/settings/receipts/create" element={<ReceiptCreatePage/>} />
          <Route path="/settings/receipts/:id/edit" element={<ReceiptCreatePage />} />
          <Route path="/settings/device-settings" element={<DeviceSettingsPage />} />
          <Route path="/settings/device-settings/create" element={<DeviceSettingCreatePage/>} />
          <Route path="/settings/device-settings/:id/edit" element={<DeviceSettingCreatePage />} />
        </Route>
        
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<DashboardPage/>} />
        

        <Route path="/calibrate" element={<CalibrationPage />} />

        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/loader" element={<SproutLoader />} />

        <Route path="/choice" element={<ChoicePage />} />
        <Route path="/tasks" element={<AppliedTasksPage />} />
        <Route path="/tasks/:id" element={<TaskDetails />} />
        <Route path="*" element={<NotFoundPage/>} />
      </Route>
    </Routes>
  );
};