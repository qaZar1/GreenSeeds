import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import { api } from "../../../api/apiProvider";
import { usePageHeader } from "../../../context/HeaderContext";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const DeviceSettingCreatePage: React.FC = () => {

  const { id } = useParams();
  const navigate = useNavigate();

  usePageHeader(
    id ? "Редактирование настройки" : "Создание настройки",
    id ? "Изменение параметра устройства" : "Добавление параметра устройства"
  );

  const [form, setForm] = useState({
    key: "",
    value: ""
  });

  const [errors, setErrors] = useState<Record<string,string>>({});
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);

  const loadSetting = async () => {
    if (!id) return;

    setLoading(true);
    setError(false);

    try {
      const data = await api.getOne("deviceSettings", id);

      setForm({
        key: data?.data.key || "",
        value: data?.data.value || ""
      });
    } catch {
      toast.error("Ошибка загрузки настройки");
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSetting();
  }, [id]);

  const update = (name: string, value: string) => {
    setForm(prev => ({
      ...prev,
      [name]: value
    }));

  };

  const validate = () => {
    const newErrors: Record<string,string> = {};

    if (!form.key.trim()) newErrors.key = "Обязательное поле";
    if (!form.value.trim()) newErrors.value = "Обязательное поле";

    setErrors(newErrors);

    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    const loading = toast.loading(
      id ? "Сохранение настройки" : "Создание настройки"
    );

    try {
      if (id) {
        await api.update("deviceSettings", {
          id: id,
          ...form
        });
      } else {
        await api.create("deviceSettings", {
          ...form
        });

      }

      toast.success(
        id ? "Настройка обновлена" : "Настройка создана",
        { id: loading }
      );

      navigate("/settings/device-settings");

    } catch {

      toast.error("Ошибка сохранения", { id: loading });

    }

  };

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadSetting} />

  return (
    <div className="w-full space-y-[24px]">
      <div className="bg-[var(--bg-card)] border border-[var(--border-color)] rounded-[16px] shadow-sm">
        <form onSubmit={handleSubmit} className="p-[24px] space-y-[20px]">

          {/* key */}
          <div>
            <label className="block text-[14px] font-medium mb-[8px] text-[var(--text-primary)]">
              Ключ *
            </label>

            <input
              type="text"
              value={form.key}
              disabled={!!id}
              onChange={e => update("key", e.target.value)}
              className={`
                w-full px-[14px] py-[10px] rounded-[10px] border
                ${errors.key ? "border-red-500" : "border-[var(--border-color)]"}
                ${id ? "bg-[var(--bg-disabled)] text-[var(--text-secondary)] cursor-not-allowed opacity-70" : "bg-[var(--bg-page)]"}
                text-[var(--text-primary)]
              `}
            />

            {errors.key && (
              <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                {errors.key}
              </p>
            )}

          </div>

          {/* value */}
          <div>
            <label className="block text-[14px] font-medium mb-[8px] text-[var(--text-primary)]">
              Значение *
            </label>
            <textarea
              value={form.value}
              onChange={e => update("value", e.target.value)}
              rows={6}
              className={`
                w-full px-[14px] py-[10px] rounded-[10px] border resize-none
                ${errors.value ? "border-red-500" : "border-[var(--border-color)]"}
                bg-[var(--bg-page)] text-[var(--text-primary)] font-mono
              `}
            />

            {errors.value && (
              <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                {errors.value}
              </p>
            )}
          </div>

          {/* buttons */}
          <div className="flex gap-[12px] pt-[8px]">
            <button
              type="button"
              onClick={() => navigate(-1)}
              className="flex-1 px-[16px] py-[10px] rounded-[10px] border border-[var(--border-color)] text-[var(--text-primary)]"
            >
              Отмена
            </button>
            <button
              type="submit"
              className="flex-1 px-[16px] py-[10px] rounded-[10px] bg-[var(--color-primary)] text-[var(--text-primary)]"
            >
              {id ? "Сохранить" : "Создать"}
            </button>
          </div>
        </form>
      </div>
    </div>

  );

};

export default DeviceSettingCreatePage;