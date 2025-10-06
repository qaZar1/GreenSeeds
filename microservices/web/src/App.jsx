import React from "react";
import './App.css'
import { Admin, Resource } from 'react-admin'
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
import UserEdit from './components/Users/UserEdit'
import UserCreate from './components/Users/UserCreate'

function App() {

  return (
    <>
      <Admin dataProvider={dataProvider} authProvider={authProvider}>
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
          edit={UserEdit}
          create={UserCreate}
          icon={GrassIcon}
          options={{ label: "Пользователи" }} 
        />
      </Admin>
    </>
  )
}

export default App
