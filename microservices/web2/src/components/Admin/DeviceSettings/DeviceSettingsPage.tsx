import React, { use, useEffect, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";

import { Table } from "../../utils/Table";

import type { Column } from "../../../types/table";
import type { Setting } from "../../../types/device-settings";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

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

      {/* add */}

      <div className="flex justify-end">

        <button
          onClick={() => navigate("/settings/device-settings/create")}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px]"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>

      </div>

      {/* table */}
      <Table
        data={settings}
        columns={columns}
        emptyMessage="Настройки ещё не созданы"
      />
    </div>
  );
};

export default DeviceSettingsPage;