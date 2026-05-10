import React from 'react';
import { useNavigate } from 'react-router-dom';

import { useAuth } from '../../../context/AuthContext';

interface LogoutActionProps {
  variant?: 'default' | 'sidebar';
}

const LogoutAction: React.FC<LogoutActionProps> = ({
  variant = 'default',
}) => {
  const navigate = useNavigate();

  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <button
      onClick={handleLogout}
      className="
        w-full flex items-center gap-[12px]
        px-[15px] py-[14px]
        text-left
        hover:bg-[var(--status-danger-bg)]
        transition-colors
        rounded-[10px]
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
        <span
          className={`
            block text-[14px] font-medium
            ${
              variant === 'sidebar'
                ? 'text-[#ffb3b3]'
                : 'text-[var(--status-danger-text)]'
            }
          `}
        >
          Выйти
        </span>

        <span
          className={`
            block text-[11px]
            ${
              variant === 'sidebar'
                ? 'text-white/60'
                : 'text-[var(--text-secondary)]'
            }
          `}
        >
          Завершить сессию
        </span>
      </div>
    </button>
  );
};

export default LogoutAction;