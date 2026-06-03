import React from "react";

interface ActionButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  icon?: string;
  disabled?: boolean;
  className?: string;
  type?: "button" | "submit";
}

const ActionButton: React.FC<ActionButtonProps> = ({
  children,
  onClick,
  icon,
  disabled = false,
  className = "",
  type = "button",
}) => {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`
        inline-flex items-center justify-center gap-[8px]

        w-full sm:w-auto

        px-[20px] py-[10px]

        bg-[var(--color-primary)]
        text-[var(--text-inverse)]

        rounded-[10px]
        font-medium
        shadow-sm

        transition-colors
        hover:bg-[var(--color-primary-hover)]

        disabled:opacity-50
        disabled:cursor-not-allowed

        ${className}
      `}
    >
      {icon && (
        <i className={`${icon} text-[14px]`} />
      )}

      {children}
    </button>
  );
};

export default ActionButton;