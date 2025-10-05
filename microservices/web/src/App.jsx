import React from "react";
import './App.css'
import { Admin, Resource } from 'react-admin'
import dataProvider from './dataProvider'
import { authProvider } from './authProvider'
import BunkerList from './components/Bunkers/Bunker'
import BunkerEdit from './components/Bunkers/BunkerEdit'
import CreateBunker from './components/Bunkers/BunkerCreate'
import WarehouseIcon from '@mui/icons-material/Warehouse';

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
      </Admin>
    </>
  )
}

export default App
