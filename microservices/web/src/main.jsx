import React from "react";
import ReactDOM from "react-dom/client";
import { Admin, Resource, ListGuesser } from "react-admin";
import jsonServerProvider from "ra-data-json-server";
import { authProvider } from "./authProvider";

const dataProvider = jsonServerProvider("https://jsonplaceholder.typicode.com");

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <Admin dataProvider={dataProvider} authProvider={authProvider}>
      <Resource name="users" list={ListGuesser} />
      <Resource name="posts" list={ListGuesser} />
    </Admin>
  </React.StrictMode>
);
