import React from "react";
import './App.css'
import { Admin, AppBar, Resource, CustomRoutes } from 'react-admin'
import dataProvider from './dataProvider'
import { authProvider } from './authProvider'
import BunkerList from './components/Admin/Bunkers/Bunker'
import BunkerEdit from './components/Admin/Bunkers/BunkerEdit'
import CreateBunker from './components/Admin/Bunkers/BunkerCreate'
import SeedList from './components/Admin/Seeds/Seed'
import SeedEdit from './components/Admin/Seeds/SeedEdit'
import CreateSeed from './components/Admin/Seeds/SeedCreate'
import UserList from './components/Admin/Users/User'
import UserCreate from './components/Admin/Users/UserCreate'
import CustomLayout from './components/AppBar/AppBar';
import { Route } from "react-router-dom";
import ProfilePage from './components/Admin/Profile/Profile';
import PlacementList from './components/Admin/Placements/Placements';
import PlacementEdit from './components/Admin/Placements/PlacementEdit';
import PlacementCreate from './components/Admin/Placements/PlacementCreate';
import ReceiptList from './components/Admin/Receipts/Receipts';
import ReceiptEdit from './components/Admin/Receipts/ReceiptEdit';
import ReceiptCreate from './components/Admin/Receipts/ReceiptCreate';
import ShiftList from './components/Admin/Shifts/Shifts';
import ShiftEdit from './components/Admin/Shifts/ShiftEdit';
import ShiftCreate from './components/Admin/Shifts/ShiftCreate';
import AssignmentsList from './components/Admin/Assignments/Assign';
import AssignmentsEdit from './components/Admin/Assignments/AssignEdit';
import AssignmentsCreate from './components/Admin/Assignments/AssignCreate';
import ReportsList from './components/Admin/Reports/Reports';
import ReportsShow from './components/Admin/Reports/ReportShow';
import ChoiceList from "./components/Operator/ChoiceTasks/Choice";
import TaskDetails from "./components/Operator/TaskDetail/TaskDetails";
import AppliedTaskList from "./components/Operator/AppliedTasks/Applied";
import { getRole } from "./authProvider";

function App() {
  const role = getRole();
  return (
    <>
      <Admin dataProvider={dataProvider} authProvider={authProvider} layout={CustomLayout}>
        {role === 'admin' && (
          <>
            <Resource
              name="bunkers"
              list={BunkerList}
              edit={BunkerEdit}
              create={CreateBunker}
            />
            <Resource
              name="seeds"
              list={SeedList}
              edit={SeedEdit}
              create={CreateSeed}
            />
            <Resource
              name="users"
              list={UserList}
              create={UserCreate}
            />
            <Resource
              name="placements"
              list={PlacementList}
              create={PlacementCreate}
              edit={PlacementEdit}
            />
            <Resource
              name="receipts"
              list={ReceiptList}
              create={ReceiptCreate}
              edit={ReceiptEdit}
            />
            <Resource
              name="shifts"
              list={ShiftList}
              create={ShiftCreate}
              edit={ShiftEdit}
            />
            <Resource
              name="assignments"
              list={AssignmentsList}
              create={AssignmentsCreate}
              edit={AssignmentsEdit}
            />
            <Resource
              name="reports"
              list={ReportsList}
              show={ReportsShow}
            />
          </>
        )}

        {role === 'operator' && (
          <>
            <Resource
              name="choice"
              list={ChoiceList}
            />

            <Resource
              name="tasks"
              list={AppliedTaskList}
            />
          </>
        )}

        <CustomRoutes>
          <Route path="/profile" element={<ProfilePage />} />
          {role === 'operator' && (
            <Route path="/tasks/:id" element={<TaskDetails />} />
          )}
        </CustomRoutes>
      </Admin>
    </>
  )
}

export default App
