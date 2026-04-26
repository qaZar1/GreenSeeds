import React, { useState, useEffect } from "react";
import type { User } from "../../../types/user";
import type { UserInput } from "../../../types/user";

interface UserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (user: UserInput) => void;
  editUser?: User | null;
}

const UserModal: React.FC<UserModalProps> = ({
  isOpen,
  onClose,
  onSave,
  editUser,
}) => {
  const [formData, setFormData] = useState<UserInput>({
    username: "",
    full_name: "",
    is_admin: false,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (editUser) {
      setFormData({
        username: editUser.username,
        full_name: editUser.full_name,
        is_admin: editUser.is_admin,
      });
    } else {
      setFormData({
        username: "",
        full_name: "",
        is_admin: false,
      });
    }

    setErrors({});
  }, [editUser, isOpen]);

  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.username.trim()) {
      newErrors.username = "Введите логин";
    }

    if (!formData.full_name.trim()) {
      newErrors.full_name = "Введите ФИО";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    onSave(formData);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-[20px]">
      <div
        className="absolute inset-0 bg-black/40 backdrop-blur-[2px]"
        onClick={onClose}
      />

      <div className="relative w-full max-w-[480px] bg-[var(--bg-card)] rounded-[16px] shadow-2xl border border-[var(--border-color)] animate-in fade-in zoom-in-95 duration-200">
        {/* Header */}
        <div className="flex items-center justify-between px-[24px] py-[20px] border-b border-[var(--border-color)]">
          <h3 className="text-[18px] font-semibold text-[var(--text-primary)]">
            {editUser ? "Редактировать пользователя" : "Добавить пользователя"}
          </h3>

          <button
            onClick={onClose}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] transition-colors"
          >
            <i className="fa-solid fa-xmark text-[16px]" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-[24px] space-y-[20px]">
          {/* Username */}
          <div>
            <label className="block text-[14px] font-medium text-[var(--text-primary)] mb-[8px]">
              Логин *
            </label>

            <input
              type="text"
              value={formData.username}
              onChange={(e) =>
                setFormData({ ...formData, username: e.target.value })
              }
              disabled={!!editUser}
              className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] ${
                errors.username
                  ? "border-[var(--status-danger-text)]"
                  : "border-[var(--border-color)]"
              }`}
              placeholder="ivanov"
            />

            {errors.username && (
              <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                {errors.username}
              </p>
            )}
          </div>

          {/* Full name */}
          <div>
            <label className="block text-[14px] font-medium text-[var(--text-primary)] mb-[8px]">
              ФИО *
            </label>

            <input
              type="text"
              value={formData.full_name}
              onChange={(e) =>
                setFormData({ ...formData, full_name: e.target.value })
              }
              className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] ${
                errors.full_name
                  ? "border-[var(--status-danger-text)]"
                  : "border-[var(--border-color)]"
              }`}
              placeholder="Иванов Иван Иванович"
            />

            {errors.full_name && (
              <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                {errors.full_name}
              </p>
            )}
          </div>

          {/* Role */}
          <div>
            <label className="block text-[14px] font-medium text-[var(--text-primary)] mb-[8px]">
              Роль
            </label>

            <div className="relative">
              <select
                value={formData.is_admin ? "admin" : "operator"}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    is_admin: e.target.value === "admin",
                  })
                }
                className="w-full px-[14px] pr-[40px] py-[10px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)] appearance-none focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
              >
                <option value="operator">Оператор</option>
                <option value="admin">Администратор</option>
              </select>

              <i className="fa-solid fa-chevron-down absolute right-[14px] top-1/2 -translate-y-1/2 text-[12px] text-[var(--text-secondary)] pointer-events-none"></i>
            </div>
          </div>

          {/* Buttons */}
          <div className="flex gap-[12px] pt-[8px]">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-[16px] py-[10px] rounded-[10px] border border-[var(--border-color)] text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors font-medium"
            >
              Отмена
            </button>

            <button
              type="submit"
              className="flex-1 px-[16px] py-[10px] rounded-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] hover:bg-[var(--color-primary-hover)] transition-colors font-medium"
            >
              {editUser ? "Сохранить" : "Добавить"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default UserModal;