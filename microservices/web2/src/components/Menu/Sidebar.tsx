import React, { useState, useRef, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const menuItems = [
  { path: '/dashboard', label: 'Дашборд', icon: 'fa-solid fa-chart-line', roles: ['admin', 'operator'] },
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
      { path: '/settings/receipts', label: 'Рецепты', icon: 'fa-solid fa-file-contract' },
      { path: '/settings/device-settings', label: 'Настройки устройства', icon: 'fa-solid fa-sliders' },
    ]
  },

  { path: '/calibrate', label: 'Калибровка', icon: 'fa-solid fa-person-walking-arrow-right', roles: ['admin'] },
];

const Sidebar = () => {
  const [openSettings, setOpenSettings] = useState(false);
  const [lineHeight, setLineHeight] = useState(0);

  const containerRef = useRef<HTMLDivElement>(null);

  const { role } = useAuth();

  if (!role) {
    throw new Error("User ID is required");
  }

  const filteredItems = menuItems.filter(item =>
    !item.roles || item.roles.includes(role)
  );

  // 👉 считаем высоту линии до последнего элемента
  useEffect(() => {
    if (openSettings && containerRef.current) {
      const lastChild = containerRef.current.querySelector('li:last-child');

      if (lastChild) {
        const rect = lastChild.getBoundingClientRect();
        const parentRect = containerRef.current.getBoundingClientRect();

        setLineHeight(rect.bottom - parentRect.top - 8); // небольшой отступ
      }
    } else {
      setLineHeight(0);
    }
  }, [openSettings]);

  return (
    <nav className="w-[250px] h-screen bg-[#34495e] text-white flex flex-col p-5 fixed left-0 top-0 overflow-y-auto">

      {/* Логотип */}
      <div className="flex items-center justify-center gap-[10px] text-[38px] font-bold mb-[30px]">
        <i className="fa-solid fa-seedling text-[#2ecc71]"></i>
        <span>Hortus</span>
      </div>

      <ul className="list-none">

        {filteredItems.map((item, index) => (
          <li key={index} className="mb-[15px]">

            {item.isSection ? (
              <>
                <button
                  onClick={() => setOpenSettings(!openSettings)}
                  className="w-full text-[#bdc3c7] flex items-center gap-[12px] p-[10px] rounded-[8px] transition-all duration-300 hover:bg-white/10 hover:text-white"
                >
                  <i className={item.icon}></i>
                  <span className="flex-1 text-left">{item.label}</span>
                  <i className={`fa-solid ${openSettings ? 'fa-chevron-up' : 'fa-chevron-down'} text-[12px]`} />
                </button>

                {/* 🔥 АНИМИРОВАННОЕ ПОДМЕНЮ */}
                <div
                  className={`
                    transition-all duration-500 overflow-hidden
                    ${openSettings ? "max-h-[500px] opacity-100 mt-[10px]" : "max-h-0 opacity-0"}
                  `}
                >
                  <div className="relative ml-[10px]" ref={containerRef}>

                    {/* линия */}
                    <div
                      className="absolute left-[6px] w-[2px] bg-white/20 transition-all duration-500"
                      style={{ height: lineHeight }}
                    />

                    <ul className="list-none pl-[20px] space-y-[10px]">
                      {item.children.map((child, childIndex) => (
                        <li key={childIndex} className="relative">

                          {/* уголок */}
                          <div className={`
                            absolute left-[-14px] top-[12px] w-[10px] h-[10px]
                            border-l-2 border-b-2 border-white/30 rounded-bl-md
                            transition-all duration-500
                            ${openSettings ? "opacity-100 scale-100" : "opacity-0 scale-75"}
                          `} />

                          <NavLink
                            to={child.path}
                            className={({ isActive }) =>
                              `no-underline text-[#bdc3c7] flex items-center gap-[12px] p-[8px] rounded-[8px] transition-all duration-300 text-[14px]
                              ${isActive ? 'bg-white/10 text-white' : 'hover:bg-white/10 hover:text-white'}`
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
                  className={({ isActive }) =>
                    `no-underline text-[#bdc3c7] flex items-center gap-[12px] p-[10px] rounded-[8px] transition-all duration-300
                    ${isActive ? 'bg-white/10 text-white' : 'hover:bg-white/10 hover:text-white'}`
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
    </nav>
  );
};

export default Sidebar;