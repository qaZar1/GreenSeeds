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
        console.log(data)
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
        <div className="flex items-center gap-[12px]">
          <div className="w-[40px] h-[40px] rounded-full bg-[var(--color-primary)] flex items-center justify-center text-[var(--text-inverse)] font-bold text-[14px]">
            {user.full_name
              .split(" ")
              .map(w => w[0])
              .join("")
              .toUpperCase()
              .slice(0, 2)}
          </div>

          <div>
            <div className="text-[14px] font-medium text-[var(--text-primary)]">
              {user.full_name}
            </div>
            <div className="text-[12px] text-[var(--text-secondary)]">
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
        <button
          onClick={handleAddUser}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px] font-medium hover:bg-[var(--color-primary-hover)] transition-colors shadow-sm"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 gap-[12px]">
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

      <Table data={users} columns={columns} emptyMessage="Пользователи не найдены" />

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