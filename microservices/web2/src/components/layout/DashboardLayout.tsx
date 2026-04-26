import { Outlet } from "react-router-dom";
import Sidebar from "../Menu/Sidebar";
import Header from "../Menu/Header";

const AppLayout = () => {
  return (
    <div id="app-layout" className="flex min-h-screen bg-[var(--bg-page)]">

      {/* Sidebar */}
      <Sidebar />

      {/* Правая часть */}
      <div className="ml-[250px] flex-1 flex flex-col">

        {/* Header */}
        <div className="p-[30px] pb-0">
          <Header/>
        </div>

        {/* Контент */}
        <main className="flex-1 p-[30px] pt-[20px]">
          <Outlet />
        </main>

      </div>
    </div>
  );
};

export default AppLayout;