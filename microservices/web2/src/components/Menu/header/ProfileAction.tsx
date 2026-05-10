import React from 'react';
import { useNavigate } from 'react-router-dom';

interface ProfileActionProps {
  variant?: 'default' | 'sidebar';
}

const ProfileAction: React.FC<ProfileActionProps> = ({
  variant = 'default',
}) => {
  const navigate = useNavigate();

  return (
    <button
      onClick={() => navigate('/profile')}
      className="
        w-full flex items-center gap-[12px]
        px-[15px] py-[14px]
        text-left
        hover:bg-[var(--bg-hover)]
        transition-colors
        rounded-[10px]
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
        <span
          className={`
            block text-[14px] font-medium
            ${
              variant === 'sidebar'
                ? 'text-white'
                : 'text-[var(--text-primary)]'
            }
          `}
        >
          Настройки
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
          Параметры профиля
        </span>
      </div>
    </button>
  );
};

export default ProfileAction;