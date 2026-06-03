import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useHeader } from '../../../context/HeaderContext';
import UserProfileButton from './UserProfileButton';
import UserProfileMenu from './UserProfileMenu';

interface HeaderProps {
  onOpenSidebar: () => void;
}

const Header: React.FC<HeaderProps> = ({ onOpenSidebar }) => {
  const { headerConfig } = useHeader();
  const navigate = useNavigate();

  const [isProfileOpen, setIsProfileOpen] = useState(false);

  return (
    <>
      <header
        className="
          flex items-center justify-between
          gap-4
          relative
        "
      >
        <div className="flex items-center gap-[12px] min-w-0">

          {/* BURGER */}
          <button
            onClick={onOpenSidebar}
            className="
              md:hidden
              w-[42px]
              h-[42px]
              rounded-[10px]
              flex
              items-center
              justify-center
              bg-[var(--bg-card)]
              border border-[var(--border-color)]
              text-[var(--text-primary)]
              shadow-sm
              flex-shrink-0
            "
          >
            <i className="fa-solid fa-bars text-[18px]" />
          </button>

          {/* TITLE */}
          <div className="min-w-0 flex-1">
            <h2
              className="
                text-[20px] md:text-[24px]
                font-bold
                text-[var(--text-primary)]

                break-words
              "
            >
              {headerConfig.title}
            </h2>

            {headerConfig.subtitle && (
              <p
                className="
                  text-[13px] md:text-[14px]
                  text-[var(--text-secondary)]

                  break-words
                "
              >
                {headerConfig.subtitle}
              </p>
            )}
          </div>

        </div>

        <div className="relative flex-shrink-0 hidden md:block">
          <UserProfileButton
            isOpen={isProfileOpen}
            onClick={() => setIsProfileOpen(prev => !prev)}
          />

          <UserProfileMenu
            isOpen={isProfileOpen}
            onClose={() => setIsProfileOpen(false)}
            onNavigate={navigate}
          />
        </div>
      </header>
    </>
  );
};

export default Header;