import React, { useEffect, useState } from "react";
import type { User } from "../../../types/user";
import DeleteConfirm from "./DeleteConfirm";
import { usePageHeader } from "../../../context/HeaderContext";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { api } from "../../../api/apiProvider";
import { useAuth } from "../../../context/AuthContext";
import toast from "react-hot-toast";
import { StatCard } from "../../utils/Card";
import FormModal from "../../utils/FormModal";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";
import ResponsiveTable from "../../utils/ResponsiveTable";
import ActionButton from "../../utils/AсtionButton";

const UsersPage: React.FC = () => {
  usePageHeader("Пользователи", "Управление доступом");

  const user = useAuth();

  const [users, setUsers] = useState<User[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const currentUsername = user?.username;

  const loadUsers = async () => {
    try {
      setLoading(true);
      setError(false);
      const data = await api.getList("users");
      if (data) {
        setUsers(data || []);
        setLoading(false);
      }
    } catch (e) {
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadUsers();
  }, []);

  const handleAddUser = () => {
    setEditingUser(null);
    setIsModalOpen(true);
  };

  const handleSaveUser = async (userData: Partial<User>) => {
    const loading = toast.loading("Сохранение пользователя...");

    try {
      await api.create("users", {
        username: userData.username,
        full_name: userData.full_name,
        is_admin: userData.is_admin || false,
      });
      setUsers(prev => [...prev, userData as User]);
      toast.success("Пользователь создан", { id: loading });
      setIsModalOpen(false);
    } catch {
      toast.error("Ошибка сохранения пользователя", { id: loading });
    }
  };

  const handleDeleteClick = (user: User) => {
    setSelectedUser(user);
    setIsDeleteOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!selectedUser) return;

    const loading = toast.loading("Удаление пользователя...");

    try {
      await api.delete("users", selectedUser.username);

      setUsers(prev => prev.filter(u => u.username !== selectedUser.username));

      toast.success(`Пользователь ${selectedUser.username} удалён`, {
        id: loading,
      });

      setSelectedUser(null);
      setIsDeleteOpen(false);

    } catch {
      toast.error("Ошибка удаления пользователя", { id: loading });
    }
  };

  const handleToggleAdmin = async (user: User, is_admin: boolean) => {
    const loading = toast.loading("Обновление роли...");

    try {
      await api.update("users", {
        username: user.username,
        full_name: user.full_name,
        is_admin,
      });

      setUsers(prev =>
        prev.map(u => (u.username === user.username ? { ...u, is_admin } : u))
      );

      toast.success("Роль обновлена", { id: loading });

    } catch {
      toast.error("Ошибка обновления роли", { id: loading });
    }
  };

  const handleResetPassword = async (user: User) => {
    const loading = toast.loading("Сброс пароля...");

    try {
      await api.update("changePass", {id: user.id})

      toast.success(`Пароль пользователя ${user.username} сброшен`, {
        id: loading,
      });

    } catch {
      toast.error("Ошибка сброса пароля", { id: loading });
    }
  };

  const columns: Column<User>[] = [
    {
      header: "ФИО",
      render: user => (
        <div className="flex items-center gap-[10px] min-w-0">
          <div className="w-[36px] h-[36px] rounded-full flex-shrink-0 bg-[var(--color-primary)] flex items-center justify-center text-[var(--text-inverse)] font-bold text-[14px]">
            {user.full_name
              .split(" ")
              .map(w => w[0])
              .join("")
              .toUpperCase()
              .slice(0, 2)}
          </div>

          <div className="min-w-0">
            <div className="text-[13px] sm:text-[14px] font-medium text-[var(--text-primary)] truncate">
              {user.full_name}
            </div>
            <div className="text-[11px] sm:text-[12px] text-[var(--text-secondary)] truncate">
              @{user.username}
            </div>
          </div>
        </div>
      ),
    },

    {
      header: "Роль",
      width: "30%",
      headerClassName: "text-center",
      className: "text-center",
      render: user => {
        const isCurrentUser = user.username === currentUsername;

        return (
          <button
            onClick={() => handleToggleAdmin(user, !user.is_admin)}
            disabled={isCurrentUser}
            className={`inline-flex items-center gap-[8px] px-[12px] py-[6px] rounded-[20px] text-[13px] font-medium transition-colors ${
              user.is_admin
                ? "bg-[var(--status-warning-bg)] text-[var(--status-warning-text)]"
                : "bg-[var(--status-info-bg)] text-[var(--status-info-text)]"
            } ${isCurrentUser ? "opacity-50 cursor-not-allowed" : "hover:opacity-80"}`}
          >
            <i className={`fa-solid ${user.is_admin ? "fa-shield-halved" : "fa-user"}`} />
            {user.is_admin ? "Администратор" : "Оператор"}
          </button>
        );
      },
    },

    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: user => {
        const isCurrentUser = user.username === currentUsername;

        return (
          <div className="flex items-center justify-end gap-[8px]">
            <button
              onClick={() => handleResetPassword(user)}
              className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)] transition-colors"
              title="Сбросить пароль"
            >
              <i className="fa-solid fa-key text-[14px]" />
            </button>

            <button
              onClick={() => handleDeleteClick(user)}
              disabled={isCurrentUser}
              className={`p-[8px] rounded-[8px] transition-colors ${
                isCurrentUser
                  ? "text-[var(--text-secondary)] opacity-30 cursor-not-allowed"
                  : "text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)]"
              }`}
              title={isCurrentUser ? "Нельзя удалить себя" : "Удалить пользователя"}
            >
              <i className="fa-solid fa-trash text-[14px]" />
            </button>
          </div>
        );
      },
    },
  ];

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadUsers} />;

  return (
    <div className="space-y-[24px] w-full">
      <div className="flex justify-end">
        <ActionButton
          onClick={handleAddUser}
          icon="fa-solid fa-plus"
        >
          Добавить
        </ActionButton>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-[12px]">
        <StatCard title="Всего" value={users.length} />
        <StatCard
          title="Администраторов"
          value={users.filter(u => u.is_admin).length}
          color="var(--status-warning-text)"
        />
        <StatCard
          title="Операторов"
          value={users.filter(u => !u.is_admin).length}
        />
      </div>

      <ResponsiveTable
  data={users}
  table={
    <Table
      data={users}
      columns={columns}
      emptyMessage="Пользователи не найдены"
    />
  }
  renderCard={(user) => {
    const isCurrentUser =
      user.username === currentUsername;

    return (
      <>
        {/* user */}
        <div className="flex items-center gap-[12px]">
          <div
            className="
              w-[42px] h-[42px]
              rounded-full
              bg-[var(--color-primary)]
              flex items-center justify-center
              text-white font-bold text-[14px]
              flex-shrink-0
            "
          >
            {user.full_name
              .split(' ')
              .map(w => w[0])
              .join('')
              .toUpperCase()
              .slice(0, 2)}
          </div>

          <div className="min-w-0 break-words">
            <div className="text-[14px] font-medium text-[var(--text-primary)] break-words">
              {user.full_name}
            </div>

            <div className="text-[12px] text-[var(--text-secondary)] break-all">
              @{user.username}
            </div>
          </div>
        </div>

        {/* role */}
        <button
          onClick={() => handleToggleAdmin(user, !user.is_admin)}
          disabled={isCurrentUser}
          className={`
            w-full
            inline-flex items-center justify-center gap-[8px]
            px-[12px] py-[10px]
            rounded-[10px]
            text-[13px]
            font-medium
            transition-colors
            ${
              user.is_admin
                ? 'bg-[var(--status-warning-bg)] text-[var(--status-warning-text)]'
                : 'bg-[var(--status-info-bg)] text-[var(--status-info-text)]'
            }
            ${
              isCurrentUser
                ? 'opacity-50 cursor-not-allowed'
                : 'hover:opacity-80'
            }
          `}
        >
          <i
            className={`fa-solid ${
              user.is_admin
                ? 'fa-shield-halved'
                : 'fa-user'
            }`}
          />

          {user.is_admin
            ? 'Администратор'
            : 'Оператор'}
        </button>

        {/* actions */}
        <div className="flex items-center gap-[10px]">
          <button
            onClick={() => handleResetPassword(user)}
            className="
              flex-1
              flex items-center justify-center gap-[8px]
              py-[10px]
              rounded-[10px]
              border border-[var(--border-color)]
              text-[var(--text-secondary)]
              hover:bg-[var(--status-info-bg)]
              hover:text-[var(--status-info-text)]
              transition-colors
            "
          >
            <i className="fa-solid fa-key" />
          </button>

          <button
            onClick={() => handleDeleteClick(user)}
            disabled={isCurrentUser}
            className={`
              flex-1
              flex items-center justify-center gap-[8px]
              py-[10px]
              rounded-[10px]
              border border-[var(--border-color)]
              transition-colors
              ${
                isCurrentUser
                  ? 'opacity-40 cursor-not-allowed text-[var(--text-secondary)]'
                  : 'text-[var(--text-secondary)] hover:bg-[var(--status-danger-bg)] hover:text-[var(--status-danger-text)]'
              }
            `}
          >
            <i className="fa-solid fa-trash" />
          </button>
        </div>
      </>
    );
  }}
/>

      {/* <UserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSaveUser}
        editUser={editingUser}
      /> */}

      <FormModal
        title="Добавить пользователя"
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSaveUser}
        fields={[
          {
            name: "username",
            label: "Логин",
            type: "text",
            required: true,
            placeholder: "ivanov"
          },
          {
            name: "full_name",
            label: "ФИО",
            type: "text",
            required: true,
            placeholder: "Иванов Иван Иванович"
          },
          {
            name: "is_admin",
            label: "Роль",
            type: "select",
            options: [
              { label: "Оператор", value: false },
              { label: "Администратор", value: true }
            ]
          }
        ]}
      />

      <DeleteConfirm
        isOpen={isDeleteOpen}
        onClose={() => setIsDeleteOpen(false)}
        onConfirm={handleConfirmDelete}
        userName={selectedUser?.full_name ?? ""}
      />
    </div>
  );
};

export default UsersPage;