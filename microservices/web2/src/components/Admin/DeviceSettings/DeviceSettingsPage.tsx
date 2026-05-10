import React, { use, useEffect, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";
import type { Column } from "../../../types/table";
import type { Setting } from "../../../types/device-settings";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";
import ActionButton from "../../utils/AсtionButton";
import ResponsiveTable from "../../utils/ResponsiveTable";
import { Table } from "../../utils/Table";

const DeviceSettingsPage: React.FC = () => {

  usePageHeader(
    "Настройки устройств",
    "Конфигурация параметров устройств"
  );

  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [settings, setSettings] = useState<Setting[]>([]);

  const [error, setError] = useState(false);
  /* ---------------- загрузка ---------------- */

  const loadSettings = async () => {
    setLoading(true);
    setError(false);

    try {
      const data = await api.getList("deviceSettings");
      setSettings(data || []);
    } catch {
      toast.error("Не удалось загрузить настройки");
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSettings();
  }, []);

  /* ---------------- delete ---------------- */

  const handleDelete = async (setting: Setting) => {
    const loading = toast.loading("Удаление настройки...");

    try {
      await api.delete("deviceSettings", setting.key);

      setSettings(prev =>
        prev.filter(s => s.key !== setting.key)
      );

      toast.success("Настройка удалена", { id: loading });

    } catch {
      toast.error("Ошибка удаления", { id: loading });
    }

  };

  /* ---------------- columns ---------------- */

  const columns: Column<Setting>[] = [

    {
      header: "Ключ",
      render: s => (
        <div className="text-[var(--text-primary)] text-[14px] font-mono">
          {s.key}
        </div>
      )
    },

    {
      header: "Значение",
      render: s => (
        <div className="text-[var(--text-primary)] text-[14px] break-all">
          {s.value}
        </div>
      )
    },
    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: s => (

        <div className="flex justify-end gap-[8px]">

          <button
            onClick={() => navigate(`/settings/device-settings/${s.key}/edit`)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)]"
          >
            <i className="fa-solid fa-pen-to-square text-[14px]" />
          </button>

          <button
            onClick={() => handleDelete(s)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)]"
          >
            <i className="fa-solid fa-trash text-[14px]" />
          </button>

        </div>

      )
    }

  ];

  /* ---------------- render ---------------- */

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadSettings} />

  return (

    <div className="space-y-[24px] w-full">

      {/* кнопка добавления */}
      <div className="flex justify-end">
        <ActionButton
          onClick={() => navigate("/settings/device-settings/create")}
          icon="fa-solid fa-plus"
        >
          Добавить
        </ActionButton>
      </div>

      <ResponsiveTable
        data={settings}
        table={
          <Table
            data={settings}
            columns={columns}
            emptyMessage="Настройки ещё не созданы"
          />
        }
        emptyMessage="Настройки ещё не созданы"
        renderCard={(s) => (
          <>
            {/* content */}
            <div className="space-y-[10px] text-[14px]">

              <div className="text-[var(--text-primary)] break-all">
                <span className="text-[var(--text-secondary)]">
                  Ключ:
                </span>{" "}
                <span className="font-mono">
                  {s.key}
                </span>
              </div>

              <div className="text-[var(--text-primary)] break-all">
                <span className="text-[var(--text-secondary)]">
                  Значение:
                </span>{" "}
                {s.value}
              </div>

            </div>

            {/* actions */}
            <div className="flex items-center gap-[10px]">

              <button
                onClick={() =>
                  navigate(`/settings/device-settings/${s.key}/edit`)
                }
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
                <i className="fa-solid fa-pen-to-square" />
              </button>

              <button
                onClick={() => handleDelete(s)}
                className="
                  flex-1
                  flex items-center justify-center gap-[8px]
                  py-[10px]
                  rounded-[10px]
                  border border-[var(--border-color)]
                  text-[var(--text-secondary)]
                  hover:bg-[var(--status-danger-bg)]
                  hover:text-[var(--status-danger-text)]
                  transition-colors
                "
              >
                <i className="fa-solid fa-trash" />
              </button>

            </div>
          </>
        )}
      />
    </div>
  );
};

export default DeviceSettingsPage;