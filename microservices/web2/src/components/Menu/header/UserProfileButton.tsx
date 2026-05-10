// components/header/UserProfileButton.tsx

import React from 'react';
import { useAuth } from '../../../context/AuthContext';

interface UserProfileButtonProps {
  isOpen: boolean;
  onClick: () => void;
}

const UserProfileButton: React.FC<UserProfileButtonProps> = ({
  isOpen,
  onClick,
}) => {
  const { full_name } = useAuth();

  const initials = full_name
    ?.split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  const parts = full_name?.split(' ');
  const shortName = `${parts?.[0] ?? ''} ${
    parts?.[1] ? parts?.[1][0] + '.' : ''
  }`;

  return (
    <button
      onClick={onClick}
      className="
        w-full
        flex items-center justify-between gap-[10px]
        bg-[var(--bg-card)]
        px-[15px]
        py-[10px]
        md:w-[240px]
        rounded-[12px]
        border border-[var(--border-color)]
        shadow-md hover:shadow-lg
        transition-all
        active:scale-[0.98]
      "
    >
      <div className="flex items-center gap-[10px] overflow-hidden">
        <div
          className="
            w-[35px] h-[35px]
            rounded-full
            bg-[var(--color-primary)]
            flex items-center justify-center
            text-[14px] font-bold
            text-[var(--text-inverse)]
            flex-shrink-0
          "
        >
          {initials}
        </div>

        <div className="text-left min-w-0">
          <span
            className="
              block
              text-[14px]
              font-bold
              text-[var(--text-primary)]
              truncate
            "
          >
            {shortName}
          </span>
        </div>
      </div>

      <i
        className={`
          fa-solid fa-chevron-down
          text-[10px]
          text-[var(--text-secondary)]
          transition-transform duration-200
          flex-shrink-0
          ${isOpen ? 'rotate-180' : ''}
        `}
      />
    </button>
  );
};

export default UserProfileButton;