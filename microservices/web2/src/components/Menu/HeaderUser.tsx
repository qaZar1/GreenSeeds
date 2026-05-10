import React, { useEffect, useRef } from 'react';

import ThemeAction from './header/ThemeAction';
import ProfileAction from './header/ProfileAction';
import LogoutAction from './header/LogoutAction';

interface UserProfileMenuProps {
  isOpen: boolean;
  onClose: () => void;
}

const UserProfileMenu: React.FC<UserProfileMenuProps> = ({
  isOpen,
  onClose,
}) => {
  const menuRef = useRef<HTMLDivElement>(null);

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

  return (
    <>
      {/* overlay */}
      <div className="fixed inset-0 z-40 bg-black/20 md:bg-transparent" />

      {/* menu */}
      <div
        ref={menuRef}
        className="
          absolute
          right-0
          top-[70px]
          w-[240px]

          bg-[var(--bg-card)]
          border border-[var(--border-color)]
          rounded-[16px]
          shadow-2xl
          overflow-hidden
          z-50

          animate-in fade-in slide-in-from-top-2
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

        <ThemeAction />

        <div className="mx-[15px] h-[1px] bg-[var(--border-color)]" />

        <ProfileAction />

        <div className="mx-[15px] h-[1px] bg-[var(--border-color)]" />

        <LogoutAction />
      </div>
    </>
  );
};

export default UserProfileMenu;