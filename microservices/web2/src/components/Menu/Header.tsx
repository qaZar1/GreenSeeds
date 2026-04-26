import React, { useState, useEffect } from "react"
import UserProfile from "./HeaderUser"
import { useNavigate } from "react-router-dom";
import { useHeader } from "../../context/HeaderContext";

const THEME_KEY = 'greenSeeds_theme';

const Header: React.FC = () => {
  const { headerConfig } = useHeader();
  const navigate = useNavigate();
  
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(THEME_KEY) as 'light' | 'dark' | null;
      if (saved) return saved;
    }
    return 'light';
  });

  useEffect(() => {
    localStorage.setItem(THEME_KEY, theme);
    document.documentElement.classList.toggle('dark', theme === 'dark');
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  return (
    <header className="flex justify-between items-center mb-[30px]">
      <div className="header-title">
        <h2 className="text-[24px] font-bold text-[var(--text-primary)]">{headerConfig.title}</h2>
        {headerConfig.subtitle && <p className="text-[14px] text-[var(--text-secondary)]">{headerConfig.subtitle}</p>}
      </div>
      
      
      <UserProfile
        currentTheme={theme}
        onToggleTheme={toggleTheme}
        onNavigate={navigate}
      />
    </header>
  )
}

export default Header