// components/header/UserProfileMenu.tsx

import React, { useEffect, useRef } from 'react';
import { useAuth } from '../../../context/AuthContext';
import { useTheme } from '../../../context/ThemeContext';

interface UserProfileMenuProps {
  isOpen: boolean;
  onClose: () => void;
  onNavigate?: (path: string) => void;
}

const UserProfileMenu: React.FC<UserProfileMenuProps> = ({
  isOpen,
  onClose,
  onNavigate,
}) => {
  const menuRef = useRef<HTMLDivElement>(null);

  const { logout } = useAuth();
  const { theme, toggleTheme } = useTheme();

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        menuRef.current &&
        !menuRef.current.contains(event.target as Node)
      ) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const handleTheme = () => {
    toggleTheme();
    onClose();
  };

  const handleSettings = () => {
    onNavigate?.('/profile');
    onClose();
  };

  const handleLogout = () => {
    logout();
    onNavigate?.('/login');
    onClose();
  };

  return (
    <>
      {/* overlay */}
      <div className="fixed inset-0 z-40 bg-black/20 md:bg-transparent" />

      {/* menu */}
      <div
        ref={menuRef}
        className="
          fixed
          bottom-4 left-4 right-4
          md:absolute md:right-0 md:left-auto md:top-[70px] md:bottom-auto
          md:w-[240px]

          bg-[var(--bg-card)]
          border border-[var(--border-color)]
          rounded-[16px]
          shadow-2xl
          overflow-hidden
          z-50

          animate-in fade-in slide-in-from-bottom-2
          md:slide-in-from-top-2
          duration-200
        "
      >
        {/* header */}
        <div
          className="
            px-[15px] py-[12px]
            bg-[var(--bg-hover)]
            border-b border-[var(--border-color)]
          "
        >
          <span
            className="
              text-[12px]
              font-semibold
              uppercase
              tracking-wide
              text-[var(--text-secondary)]
            "
          >
            Меню
          </span>
        </div>

        {/* theme */}
        <button
          onClick={handleTheme}
          className="
            w-full flex items-center gap-[12px]
            px-[15px] py-[14px]
            text-left
            hover:bg-[var(--bg-hover)]
            transition-colors
          "
        >
          <div
            className="
              w-[32px] h-[32px]
              rounded-[8px]
              bg-[var(--status-info-bg)]
              flex items-center justify-center
              flex-shrink-0
            "
          >
            <i
              className={`
                fa-solid
                ${
                  theme === 'dark'
                    ? 'fa-sun text-[var(--status-warning-text)]'
                    : 'fa-moon text-[var(--status-info-text)]'
                }
              `}
            />
          </div>

          <div className="min-w-0">
            <span className="block text-[14px] font-medium text-[var(--text-primary)]">
              {theme === 'dark'
                ? 'Светлая тема'
                : 'Тёмная тема'}
            </span>

            <span className="block text-[11px] text-[var(--text-secondary)]">
              Переключить интерфейс
            </span>
          </div>
        </button>

        <div className="mx-[15px] h-[1px] bg-[var(--border-color)]" />

        {/* settings */}
        <button
          onClick={handleSettings}
          className="
            w-full flex items-center gap-[12px]
            px-[15px] py-[14px]
            text-left
            hover:bg-[var(--bg-hover)]
            transition-colors
          "
        >
          <div
            className="
              w-[32px] h-[32px]
              rounded-[8px]
              bg-[var(--status-warning-bg)]
              flex items-center justify-center
              flex-shrink-0
            "
          >
            <i className="fa-solid fa-gear text-[var(--status-warning-text)]" />
          </div>

          <div>
            <span className="block text-[14px] font-medium text-[var(--text-primary)]">
              Настройки
            </span>

            <span className="block text-[11px] text-[var(--text-secondary)]">
              Параметры профиля
            </span>
          </div>
        </button>

        <div className="mx-[15px] h-[1px] bg-[var(--border-color)]" />

        {/* logout */}
        <button
          onClick={handleLogout}
          className="
            w-full flex items-center gap-[12px]
            px-[15px] py-[14px]
            text-left
            hover:bg-[var(--status-danger-bg)]
            transition-colors
          "
        >
          <div
            className="
              w-[32px] h-[32px]
              rounded-[8px]
              bg-[var(--status-danger-bg)]
              flex items-center justify-center
              flex-shrink-0
            "
          >
            <i className="fa-solid fa-right-from-bracket text-[var(--status-danger-text)]" />
          </div>

          <div>
            <span className="block text-[14px] font-medium text-[var(--status-danger-text)]">
              Выйти
            </span>

            <span className="block text-[11px] text-[var(--text-secondary)]">
              Завершить сессию
            </span>
          </div>
        </button>
      </div>
    </>
  );
};

export default UserProfileMenu;