import React from 'react';
import { useTheme } from '../../../context/ThemeContext';

interface ThemeActionProps {
  variant?: 'default' | 'sidebar';
}

const ThemeAction: React.FC<ThemeActionProps> = ({
  variant = 'default',
}) => {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
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
          {theme === 'dark'
            ? 'Светлая тема'
            : 'Тёмная тема'}
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
          Переключить интерфейс
        </span>
      </div>
    </button>
  );
};

export default ThemeAction;