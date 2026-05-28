import React, { useState, useRef } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import ProfileAction from './header/ProfileAction';
import LogoutAction from './header/LogoutAction';
import ThemeAction from './header/ThemeAction';

const menuItems = [
  { path: '/users', label: 'Пользователи', icon: 'fa-solid fa-users', roles: ['admin'] },
  { path: '/choice', label: 'Выбор задания', icon: 'fa-solid fa-hand-pointer', roles: ['operator'] },
  { path: '/tasks', label: 'Задания на смену', icon: 'fa-solid fa-clipboard-list', roles: ['operator'] },
  { path: '/shifts', label: 'План производства', icon: 'fa-solid fa-user-clock', roles: ['admin'] },
  { path: '/assignments', label: 'Сменные задания', icon: 'fa-solid fa-clipboard-check', roles: ['admin'] },
  { path: '/reports', label: 'Отчеты', icon: 'fa-solid fa-file-invoice', roles: ['admin'] },
  { path: '/logs', label: 'Логи', icon: 'fa-solid fa-list-alt', roles: ['admin'] },

  {
    label: 'Настройки',
    icon: 'fa-solid fa-gear',
    roles: ['admin'],
    isSection: true,
    children: [
      { path: '/settings/bunkers', label: 'Бункеры', icon: 'fa-solid fa-warehouse' },
      { path: '/settings/seeds', label: 'Семена', icon: 'fa-solid fa-seedling' },
      { path: '/settings/placements', label: 'Расположение', icon: 'fa-solid fa-link' },
      { path: '/settings/recipes', label: 'Рецепты', icon: 'fa-solid fa-file-contract' },
      { path: '/settings/device-settings', label: 'Настройки устройства', icon: 'fa-solid fa-sliders' },
    ]
  },

  { path: '/calibrate', label: 'Калибровка', icon: 'fa-solid fa-person-walking-arrow-right', roles: ['admin'] },
];

interface SidebarProps {
  isMobileOpen: boolean;
  setIsMobileOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const Sidebar: React.FC<SidebarProps> = ({
  isMobileOpen,
  setIsMobileOpen,
}) => {
  const [openSettings, setOpenSettings] = useState(false);

  const containerRef = useRef<HTMLDivElement>(null);

  const { role, logout } = useAuth();
  const navigate = useNavigate();

  if (!role) {
    throw new Error('User ID is required');
  }

  const filteredItems = menuItems.filter(
    item => !item.roles || item.roles.includes(role)
  );

  return (
    <>
      {/* OVERLAY */}
      {isMobileOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 md:hidden"
          onClick={() => setIsMobileOpen(false)}
        />
      )}

      {/* SIDEBAR */}
      <nav
        className={`
          fixed top-0 left-0 z-50
          w-screen md:w-[250px]
          h-screen
          bg-[#34495e]
          text-white
          flex flex-col
          p-5
          overflow-y-auto
          transition-transform duration-300

          ${
            isMobileOpen
              ? 'translate-x-0'
              : '-translate-x-full'
          }

          md:translate-x-0
        `}
      >
        {/* HEADER */}
        <div className="flex items-center justify-between mb-[30px]">
          <div className="flex items-center gap-[10px] text-[38px] font-bold">
            <i className="fa-solid fa-seedling text-[#2ecc71]"></i>
            <span>Hortus</span>
          </div>

          {/* CLOSE BUTTON MOBILE */}
          <button
            onClick={() => setIsMobileOpen(false)}
            className="
              md:hidden
              flex
              items-center
              justify-center
              w-[32px]
              h-[32px]
              text-white
              text-[22px]
              leading-none
            "
          >
            <i className="fa-solid fa-xmark"></i>
          </button>
        </div>

        <ul className="list-none">
          {filteredItems.map((item, index) => (
            <li key={index} className="mb-[15px]">

              {item.isSection ? (
                <>
                  <button
                    onClick={() => {
                      setOpenSettings(!openSettings)
                    }}
                    className="
                      w-full
                      text-[#bdc3c7]
                      flex
                      items-center
                      gap-[12px]
                      p-[10px]
                      rounded-[8px]
                      transition-all
                      duration-300
                      hover:bg-white/10
                      hover:text-white
                    "
                  >
                    <i className={item.icon}></i>

                    <span className="flex-1 text-left">
                      {item.label}
                    </span>

                    <i
                      className={`fa-solid ${
                        openSettings
                          ? 'fa-chevron-up'
                          : 'fa-chevron-down'
                      } text-[12px]`}
                    />
                  </button>

                  <div
                    className={`
                      transition-all
                      duration-500
                      overflow-hidden
                      ${
                        openSettings
                          ? 'max-h-[500px] opacity-100 mt-[10px]'
                          : 'max-h-0 opacity-0'
                      }
                    `}
                  >
                    <div className="relative ml-[10px]" ref={containerRef}>
                      <ul className="list-none pl-[20px] space-y-[10px]">

                        {item.children.map((child, childIndex) => (
                          <li
                            key={childIndex}
                            className="relative"
                          >

                            {childIndex !== item.children.length - 1 && (
                              <div className="absolute left-[-14px] top-[20px] w-[2px] h-full bg-white/20" />
                            )}

                            <div
                              className={`
                                absolute left-[-14px] top-[12px]
                                w-[10px] h-[10px]
                                border-l-2 border-b-2
                                border-white/30
                                rounded-bl-md
                                transition-all duration-500
                                ${
                                  openSettings
                                    ? 'opacity-100 scale-100'
                                    : 'opacity-0 scale-75'
                                }
                              `}
                            />

                            <NavLink
                              to={child.path}
                              onClick={() => setIsMobileOpen(false)}
                              className={({ isActive }) =>
                                `
                                no-underline
                                text-[#bdc3c7]
                                flex
                                items-center
                                gap-[12px]
                                p-[8px]
                                rounded-[8px]
                                transition-all
                                duration-300
                                text-[14px]

                                ${
                                  isActive
                                    ? 'bg-white/10 text-white'
                                    : 'hover:bg-white/10 hover:text-white'
                                }
                              `
                              }
                            >
                              <i className={child.icon}></i>
                              <span>{child.label}</span>
                            </NavLink>

                          </li>
                        ))}

                      </ul>
                    </div>
                  </div>
                </>
              ) : (
                item.path && (
                  <NavLink
                    to={item.path}
                    onClick={() => setIsMobileOpen(false)}
                    className={({ isActive }) =>
                      `
                      no-underline
                      text-[#bdc3c7]
                      flex
                      items-center
                      gap-[12px]
                      p-[10px]
                      rounded-[8px]
                      transition-all
                      duration-300

                      ${
                        isActive
                          ? 'bg-white/10 text-white'
                          : 'hover:bg-white/10 hover:text-white'
                      }
                    `
                    }
                  >
                    <i className={item.icon}></i>
                    <span>{item.label}</span>
                  </NavLink>
                )
              )}

            </li>
          ))}
        </ul>

        {/* MOBILE USER ACTIONS */}
        <div className="mt-auto pt-[20px] border-t border-white/10 md:hidden">
          {/* Theme */}
          <ThemeAction variant="sidebar"/>

          {/* Profile */}
          <ProfileAction variant="sidebar"/>

          {/* Logout */}
          <LogoutAction variant="sidebar"/>

        </div>
      </nav>
    </>
  );
};

export default Sidebar;