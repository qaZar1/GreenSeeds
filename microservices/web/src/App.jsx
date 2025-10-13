import React from "react";
import './App.css'
import { Admin, AppBar, Resource, CustomRoutes } from 'react-admin'
import dataProvider from './dataProvider'
import { authProvider } from './authProvider'
import BunkerList from './components/Bunkers/Bunker'
import BunkerEdit from './components/Bunkers/BunkerEdit'
import CreateBunker from './components/Bunkers/BunkerCreate'
import WarehouseIcon from '@mui/icons-material/Warehouse';
import SeedList from './components/Seeds/Seed'
import SeedEdit from './components/Seeds/SeedEdit'
import CreateSeed from './components/Seeds/SeedCreate'
import GrassIcon from '@mui/icons-material/Grass';
import UserList from './components/Users/User'
import UserCreate from './components/Users/UserCreate'
import PersonIcon from '@mui/icons-material/Person';
import CustomLayout from './components/AppBar/AppBar';
import { Route } from "react-router-dom";
import ProfilePage from './components/Profile/Profile';
import PlacementList from './components/Placements/Placements';
import PlacementEdit from './components/Placements/PlacementEdit';
import PlacementCreate from './components/Placements/PlacementCreate';
import LinkIcon from '@mui/icons-material/Link';
import ReceiptList from './components/Receipts/Receipts';
import ReceiptEdit from './components/Receipts/ReceiptEdit';
import ReceiptCreate from './components/Receipts/ReceiptCreate';
import ReceiptIcon from '@mui/icons-material/Receipt';

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
          options={{ label: "Связи бункер-семена" }} 
        />
        <Resource
          name="receipts"
          list={ReceiptList}
          create={ReceiptCreate}
          edit={ReceiptEdit}
          icon={ReceiptIcon}
          options={{ label: "Рецепты" }} 
        />
        <CustomRoutes>
          <Route path="/profile" element={<ProfilePage />} />
        </CustomRoutes>
      </Admin>
    </>
  )
}

export default App
