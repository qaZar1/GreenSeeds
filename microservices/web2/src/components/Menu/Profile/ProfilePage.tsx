import React, { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { api } from "../../../api/apiProvider";
import { usePageHeader } from "../../../context/HeaderContext";
import FormModal from "../../utils/FormModal";
import { useAuth } from "../../../context/AuthContext";
import SproutLoader from "../../utils/Loader/SproutLoader";

type Profile = {
  username: string;
  full_name: string;
  is_admin: boolean;
};

const ProfilePage: React.FC = () => {
  usePageHeader("Профиль", "Управление учетной записью");

  const { user_id } = useAuth();

  const [profile, setProfile] = useState<Profile | null>(null);
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false);

  const [loading, setLoading] = useState(true);

  /* ---------- загрузка профиля ---------- */

  useEffect(() => {
    if (!user_id) return;

    const loadProfile = async () => {
      try {
        const res = await api.getOne("users", user_id);
        setProfile(res.data);
        setLoading(false);
      } catch {
        toast.error("Ошибка загрузки профиля");
        setLoading(false);
      }
    };

    loadProfile();
  }, [user_id]);

  /* ---------- сохранение профиля ---------- */

  const handleSave = async () => {
    if (!profile) return;

    const loading = toast.loading("Сохранение профиля...");

    try {
      await api.update("users", {
        username: profile.username,
        full_name: profile.full_name,
        is_admin: profile.is_admin,
      });

      toast.success("Профиль обновлен", { id: loading });
    } catch {
      toast.error("Ошибка сохранения профиля", { id: loading });
    }
  };

  /* ---------- смена пароля ---------- */

  const handlePasswordChange = async (data: Record<string, any>) => {
    if (data.newPassword !== data.confirmPassword) {
      toast.error("Новые пароли не совпадают");
      return;
    }

    const loading = toast.loading("Смена пароля...");

    try {
      await api.update("changePass", {
        user_id: user_id,
        old_password: data.oldPassword,
        new_password: data.newPassword,
      })

      toast.success("Пароль обновлен", { id: loading });
      setIsPasswordModalOpen(false);
    } catch {
      toast.error("Ошибка смены пароля", { id: loading });
    }
  };

  if (loading) return <SproutLoader/>
  if (!profile) return null;

  return (
    <div className="space-y-[24px] max-w-[700px] mx-auto">

      {/* карточка профиля */}
      <div className="bg-[var(--bg-card)] border border-[var(--border-color)] rounded-[12px] p-[24px] shadow-sm">

        <div className="flex items-center justify-between mb-[20px]">

          <div className="flex items-center gap-[12px]">
            <div className="w-[40px] h-[40px] rounded-[10px] bg-[var(--color-primary)] flex items-center justify-center text-white">
              <i className="fa-solid fa-user text-[16px]" />
            </div>

            <div>
              <div className="text-[16px] font-semibold text-[var(--text-primary)]">
                Основная информация
              </div>
              <div className="text-[13px] text-[var(--text-secondary)]">
                Данные вашей учетной записи
              </div>
            </div>
          </div>

          <button
            onClick={() => setIsPasswordModalOpen(true)}
            className="inline-flex items-center gap-[8px] px-[14px] py-[8px] text-[13px] rounded-[8px] border border-[var(--border-color)] text-[var(--text-secondary)] hover:text-[var(--status-warning-text)] hover:bg-[var(--status-warning-bg)] transition-colors"
          >
            <i className="fa-solid fa-key text-[12px]" />
            Сменить пароль
          </button>

        </div>

        {/* форма */}
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleSave();
          }}
          className="space-y-[16px]"
        >

          {/* username */}
          <div>
            <label className="text-[13px] text-[var(--text-secondary)] block mb-[6px]">
              Имя пользователя
            </label>

            <input
              value={profile?.username ?? ""}
              disabled
              className="w-full px-[12px] py-[10px] rounded-[8px] border border-[var(--border-light)] bg-[var(--bg-page)] text-[var(--text-secondary)]"
            />
          </div>

          {/* full name */}
          <div>
            <label className="text-[13px] text-[var(--text-secondary)] block mb-[6px]">
              ФИО
            </label>

            <input
              value={profile?.full_name ?? ""}
              onChange={(e) =>
                setProfile(prev =>
                  prev ? { ...prev, full_name: e.target.value } : prev
                )
              }
              className="w-full px-[12px] py-[10px] rounded-[8px] border border-[var(--border-color)] bg-[var(--bg-card)] text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
            />
          </div>

          {/* admin */}
          <div className="flex items-center gap-[10px] pt-[6px]">
            <div
              className={`px-[10px] py-[4px] rounded-[6px] text-[12px] font-medium ${
                profile.is_admin
                  ? "bg-[var(--status-info-bg)] text-[var(--status-info-text)]"
                  : "bg-[var(--border-light)] text-[var(--text-secondary)]"
              }`}
            >
              {profile.is_admin ? "Администратор" : "Пользователь"}
            </div>
          </div>

          {/* кнопка */}
          <div className="pt-[10px] flex justify-end">
            <button
              type="submit"
              className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px] font-medium hover:bg-[var(--color-primary-hover)] transition-colors shadow-sm"
            >
              <i className="fa-solid fa-save text-[14px]" />
              Сохранить
            </button>
          </div>

        </form>

      </div>

      {/* модалка смены пароля */}
      <FormModal
        title="Смена пароля"
        isOpen={isPasswordModalOpen}
        onClose={() => setIsPasswordModalOpen(false)}
        onSubmit={handlePasswordChange}
        fields={[
          {
            name: "oldPassword",
            label: "Старый пароль",
            type: "password",
            required: true,
          },
          {
            name: "newPassword",
            label: "Новый пароль",
            type: "password",
            required: true,
          },
          {
            name: "confirmPassword",
            label: "Повторите пароль",
            type: "password",
            required: true,
          },
        ]}
      />

    </div>
  );
};

export default ProfilePage;