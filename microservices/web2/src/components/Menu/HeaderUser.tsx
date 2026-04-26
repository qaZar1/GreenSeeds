import React, { useState, useRef, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';

interface UserProfileProps {
  name?: string;
  email?: string;
  onLogout?: () => void;
  onToggleTheme?: () => void;
  onNavigate?: (path: string) => void;
  currentTheme?: 'light' | 'dark';
}

const UserProfile: React.FC<UserProfileProps> = ({ 
  onToggleTheme,
  onNavigate,
  currentTheme = 'light'
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const { full_name } = useAuth();
  const { logout } = useAuth();

  const initials = full_name?.split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
  
  const parts = full_name?.split(" ");
  const shortName = `${parts?.[0]} ${parts?.[1] ? parts?.[1][0] + "." : ""}`;

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleToggleTheme = () => {
    onToggleTheme?.();
    setIsOpen(false);
  };

  const handleSettings = () => {
    onNavigate?.('/profile');
    setIsOpen(false);
  };

  const handleLogout = () => {
    logout();
    setIsOpen(false);
    onNavigate?.('/login');
  };

  return (
    <div className="relative w-[240px]" ref={dropdownRef}>
      
      {/* Кнопка профиля */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex items-center justify-between gap-[10px] bg-[var(--bg-card)] px-[15px] py-[10px] rounded-[12px] shadow-md cursor-pointer hover:shadow-lg transition-all border border-[var(--border-color)] outline-none active:scale-[0.98]"
      >
        <div className="flex items-center gap-[10px] overflow-hidden">
          <div className="w-[35px] h-[35px] bg-[var(--color-primary)] rounded-full flex items-center justify-center text-[var(--text-inverse)] font-bold text-[14px] shadow-sm flex-shrink-0">
            {initials}
          </div>
          <div className="text-left min-w-0">
            <span className="block text-[14px] font-bold text-[var(--text-primary)] truncate">{shortName}</span>
          </div>
        </div>
        <i className={`fa-solid fa-chevron-down text-[10px] text-[var(--text-secondary)] transition-transform duration-200 flex-shrink-0 ${isOpen ? 'rotate-180' : ''}`}></i>
      </button>

      {/* Выпадающее меню */}
      {isOpen && (
        <>
          <div className="absolute right-0 mt-[12px] w-[240px] bg-[var(--bg-card)] rounded-[12px] shadow-xl border border-[var(--border-color)] z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200">
            
            {/* Заголовок меню */}
            <div className="px-[15px] py-[12px] bg-[var(--bg-hover)] border-b border-[var(--border-color)]">
              <span className="text-[12px] font-semibold text-[var(--text-secondary)] uppercase tracking-wide">Меню</span>
            </div>

            {/* Пункт: Смена темы */}
            <button
              onClick={handleToggleTheme}
              className="w-full flex items-center gap-[12px] px-[15px] py-[14px] text-left hover:bg-[var(--bg-hover)] transition-colors group"
            >
              <div className="w-[32px] h-[32px] rounded-[8px] bg-[var(--status-info-bg)] flex items-center justify-center group-hover:bg-[var(--bg-hover)] transition-colors flex-shrink-0">
                <i className={`fa-solid ${currentTheme === 'dark' ? 'fa-sun text-[var(--status-warning-text)]' : 'fa-moon text-[var(--status-info-text)]'} text-[16px]`}></i>
              </div>
              <div className="flex-1 min-w-0">
                <span className="block text-[14px] font-medium text-[var(--text-primary)]">
                  {currentTheme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}
                </span>
                <span className="block text-[11px] text-[var(--text-secondary)]">Переключить интерфейс</span>
              </div>
            </button>

            {/* Разделитель */}
            <div className="h-[1px] bg-gradient-to-r from-transparent via-[var(--border-color)] to-transparent mx-[15px]"></div>

            {/* Пункт: Настройки */}
            <button
              onClick={handleSettings}
              className="w-full flex items-center gap-[12px] px-[15px] py-[14px] text-left hover:bg-[var(--bg-hover)] transition-colors group"
            >
              <div className="w-[32px] h-[32px] rounded-[8px] bg-[var(--status-warning-bg)] flex items-center justify-center group-hover:bg-[var(--bg-hover)] transition-colors flex-shrink-0">
                <i className="fa-solid fa-gear text-[var(--status-warning-text)] text-[16px]"></i>
              </div>
              <div className="flex-1 min-w-0">
                <span className="block text-[14px] font-medium text-[var(--text-primary)]">Настройки</span>
                <span className="block text-[11px] text-[var(--text-secondary)]">Параметры профиля</span>
              </div>
            </button>

            {/* Разделитель */}
            <div className="h-[1px] bg-gradient-to-r from-transparent via-[var(--border-color)] to-transparent mx-[15px]"></div>

            {/* Пункт: Выйти */}
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-[12px] px-[15px] py-[14px] text-left hover:bg-[var(--status-danger-bg)] transition-colors group"
            >
              <div className="w-[32px] h-[32px] rounded-[8px] bg-[var(--status-danger-bg)] flex items-center justify-center group-hover:opacity-90 transition-colors flex-shrink-0">
                <i className="fa-solid fa-right-from-bracket text-[var(--status-danger-text)] text-[16px]"></i>
              </div>
              <div className="flex-1 min-w-0">
                <span className="block text-[14px] font-medium text-[var(--status-danger-text)]">Выйти</span>
                <span className="block text-[11px] text-[var(--text-secondary)]">Завершить сессию</span>
              </div>
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default UserProfile;