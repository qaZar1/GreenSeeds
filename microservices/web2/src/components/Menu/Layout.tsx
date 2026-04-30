import React, { useState, useEffect } from "react";
import Header from "./Header";
import Sidebar from "./Sidebar";

const THEME_KEY = 'greenSeeds_theme';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(THEME_KEY) as 'light' | 'dark' | null;
      return saved || 'light';
    }
    return 'light';
  });

  // Применяем тему
  useEffect(() => {
    localStorage.setItem(THEME_KEY, theme);
    document.documentElement.classList.toggle('dark', theme === 'dark');
  }, [theme]);

  // Блокируем скролл тела при открытом мобильном меню
  useEffect(() => {
    if (window.innerWidth < 1024) {
      document.body.style.overflow = isSidebarOpen ? "hidden" : "";
    }
    return () => { document.body.style.overflow = ""; };
  }, [isSidebarOpen]);

  const toggleTheme = () => setTheme(prev => prev === 'light' ? 'dark' : 'light');

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      {/* Сайдбар */}
      <Sidebar/>

      {/* Основной контент */}
      <div className="lg:ml-[250px] min-h-screen">
        <Header />
        <main className="p-[20px]">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;