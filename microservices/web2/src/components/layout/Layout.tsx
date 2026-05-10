import { Outlet } from "react-router-dom";
import Sidebar from "../Menu/Sidebar";
import Header from "../Menu/header/Header";
import { useState } from "react";

const AppLayout = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  return (
    <div id="app-layout" className="flex min-h-screen bg-[var(--bg-page)]">

      {/* Sidebar */}
      <Sidebar
        isMobileOpen={isSidebarOpen}
        setIsMobileOpen={setIsSidebarOpen}
      />

      {/* Правая часть */}
      <div className="ml-0 md:ml-[250px] flex-1 flex flex-col min-w-0">

        {/* Header */}
        <div className="p-[16px] md:p-[30px] pt-[20px] pb-0">
          <Header onOpenSidebar={() => setIsSidebarOpen(true)} />
        </div>

        {/* Контент */}
        <main className="flex-1 p-[16px] md:p-[30px] pt-[20px] min-w-0">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default AppLayout;