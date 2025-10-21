import React from "react";
import './App.css'
import { Admin, AppBar, Resource, CustomRoutes } from 'react-admin'
import dataProvider from './dataProvider'
import { authProvider } from './authProvider'
import BunkerList from './components/Admin/Bunkers/Bunker'
import BunkerEdit from './components/Admin/Bunkers/BunkerEdit'
import CreateBunker from './components/Admin/Bunkers/BunkerCreate'
import WarehouseIcon from '@mui/icons-material/Warehouse';
import SeedList from './components/Admin/Seeds/Seed'
import SeedEdit from './components/Admin/Seeds/SeedEdit'
import CreateSeed from './components/Admin/Seeds/SeedCreate'
import GrassIcon from '@mui/icons-material/Grass';
import UserList from './components/Admin/Users/User'
import UserCreate from './components/Admin/Users/UserCreate'
import PersonIcon from '@mui/icons-material/Person';
import CustomLayout from './components/AppBar/AppBar';
import { Route } from "react-router-dom";
import ProfilePage from './components/Admin/Profile/Profile';
import PlacementList from './components/Admin/Placements/Placements';
import PlacementEdit from './components/Admin/Placements/PlacementEdit';
import PlacementCreate from './components/Admin/Placements/PlacementCreate';
import LinkIcon from '@mui/icons-material/Link';
import ReceiptList from './components/Admin/Receipts/Receipts';
import ReceiptEdit from './components/Admin/Receipts/ReceiptEdit';
import ReceiptCreate from './components/Admin/Receipts/ReceiptCreate';
import ReceiptIcon from '@mui/icons-material/Receipt';
import ShiftList from './components/Admin/Shifts/Shifts';
import ShiftEdit from './components/Admin/Shifts/ShiftEdit';
import ShiftCreate from './components/Admin/Shifts/ShiftCreate';
import DragIndicatorIcon from '@mui/icons-material/DragIndicator';
import AssignmentsList from './components/Admin/Assignments/Assign';
import AssignmentsEdit from './components/Admin/Assignments/AssignEdit';
import AssignmentsCreate from './components/Admin/Assignments/AssignCreate';
import AssignmentIcon from '@mui/icons-material/Assignment';
import ReportsList from './components/Admin/Reports/Reports';
import ReportsShow from './components/Admin/Reports/ReportShow';
import PageWithTasks from "./components/Operator/ChoiceTasks/PageWithTasks";

function App() {
  return (
    <>
      <Admin dataProvider={dataProvider} authProvider={authProvider} layout={CustomLayout}>
        <Resource
          name="bunkers"
          list={BunkerList}
          edit={BunkerEdit}
          create={CreateBunker}
          icon={WarehouseIcon}
          options={{ label: "Бункеры" }} 
        />
        <Resource
          name="seeds"
          list={SeedList}
          edit={SeedEdit}
          create={CreateSeed}
          icon={GrassIcon}
          options={{ label: "Семена" }} 
        />
        <Resource
          name="users"
          list={UserList}
          create={UserCreate}
          icon={PersonIcon}
          options={{ label: "Пользователи" }} 
        />
        <Resource
          name="placements"
          list={PlacementList}
          create={PlacementCreate}
          edit={PlacementEdit}
          icon={LinkIcon}
          options={{ label: "Расположение семян" }} 
        />
        <Resource
          name="receipts"
          list={ReceiptList}
          create={ReceiptCreate}
          edit={ReceiptEdit}
          icon={ReceiptIcon}
          options={{ label: "Рецепты" }} 
        />
        <Resource
          name="shifts"
          list={ShiftList}
          create={ShiftCreate}
          edit={ShiftEdit}
          icon={DragIndicatorIcon}
          options={{ label: "Смены" }} 
        />
        <Resource
          name="assignments"
          list={AssignmentsList}
          create={AssignmentsCreate}
          edit={AssignmentsEdit}
          icon={AssignmentIcon}
          options={{ label: "Задания" }} 
        />
        <Resource
          name="reports"
          list={ReportsList}
          show={ReportsShow}
          icon={AssignmentIcon}
          options={{ label: "Отчеты" }} 
        />


        <Resource
          name="tasks"
          list={PageWithTasks}
          options={{ label: "Задания" }} 
        />

        <CustomRoutes>
          <Route path="/profile" element={<ProfilePage />} />
        </CustomRoutes>
      </Admin>
    </>
  )
}

export default App
