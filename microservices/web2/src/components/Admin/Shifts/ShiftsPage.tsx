import React, { useEffect, useState } from "react";
import type { Shift } from "../../../types/shifts";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import { StatCard } from "../../utils/Card";
import FormModal from "../../utils/FormModal";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const ShiftsPage: React.FC = () => {
  usePageHeader("План производства", "Планируемые смены на сегодня и будущие даты");

  const [shifts, setShifts] = useState<Shift[]>([]);
  const [editingShift, setEditingShift] = useState<Shift | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [showHistory, setShowHistory] = useState(false);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadShifts = async () => {
    setLoading(true);
    setError(false);

    try {
      const data = await api.getList("shifts");
      setShifts(data || []);
    } catch {
      setError(true);
      toast.error("Ошибка загрузки смен");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadShifts();
  }, []);

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const upcomingShifts = shifts
    .filter(s => new Date(s.dt) >= today)
    .sort((a, b) => new Date(a.dt).getTime() - new Date(b.dt).getTime());

  const historyShifts = shifts
    .filter(s => new Date(s.dt) < today)
    .sort((a, b) => new Date(b.dt).getTime() - new Date(a.dt).getTime());

  const dt = (date: string) => {
    return new Date(date).toLocaleString("ru-RU");
  };

  const handleAddShift = () => {
    setEditingShift(null);
    setIsModalOpen(true);
  };

  const handleEdit = (shift: Shift) => {
    setEditingShift(shift);
    setIsModalOpen(true);
  };

  const handleSave = async (data: Record<string, any>) => {
    const dt = data.dt;

    const loading = toast.loading(
      editingShift ? "Изменение смены" : "Добавление смены"
    );

    try {
      if (editingShift) {
        await api.update("shifts", {
          shift: editingShift.shift,
          dt: dt,
        });

        toast.success("Смена обновлена", { id: loading });
      } else {
        await api.create("shifts", {
          dt: dt,
        });

        toast.success("Смена создана", { id: loading });
      }

      await loadShifts();

      setEditingShift(null);
      setIsModalOpen(false);
    } catch {
      toast.error("Ошибка сохранения смены", { id: loading });
    }
  };

  const handleDelete = async (shift: Shift) => {
    const loading = toast.loading("Удаление смены...");

    try {
      await api.delete("shifts", shift.shift);

      setShifts(prev => prev.filter(s => s.shift !== shift.shift));

      toast.success("Смена удалена", { id: loading });
    } catch {
      toast.error("Ошибка удаления смены", { id: loading });
    }
  };

  const columns: Column<Shift>[] = [
    {
      header: "Дата",
      render: shift => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {dt(shift.dt)}
        </div>
      ),
    },
    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: shift => {
        const isPast = new Date(shift.dt) < today;

        return (
          <div className="flex items-center justify-end gap-[8px]">

            <button
              onClick={() => handleEdit(shift)}
              disabled={isPast}
              className={`p-[8px] rounded-[8px] transition-colors
                ${isPast
                  ? "text-gray-400 cursor-not-allowed"
                  : "text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)]"
                }`}
              title={isPast ? "Нельзя редактировать прошедшую смену" : "Изменить время"}
            >
              <i className="fa-solid fa-clock text-[14px]" />
            </button>

            <button
              onClick={() => handleDelete(shift)}
              className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)] transition-colors"
              title="Удалить смену"
            >
              <i className="fa-solid fa-trash text-[14px]" />
            </button>

          </div>
        );
      }
    },
  ];

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadShifts} />;

  return (
    <div className="space-y-[24px] w-full">

      {/* кнопка добавления */}
      <div className="flex justify-end">
        <button
          onClick={handleAddShift}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px] font-medium hover:bg-[var(--color-primary-hover)] transition-colors shadow-sm"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>
      </div>

      {/* статистика */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-[12px]">
        <StatCard title="Всего смен" value={shifts.length} />
        <StatCard
          title="Ближайшие смены"
          value={upcomingShifts.length}
          color="var(--status-success-text)"
        />
      </div>

      {/* ближайшие смены */}
      <div>
        <Table
          data={upcomingShifts}
          columns={columns}
          emptyMessage="Ближайшие смены отсутствуют"
        />
      </div>

      {/* история */}
      {historyShifts.length > 0 && (
        <div className="border-t border-[var(--border-color)] pt-[16px]">

          <button
            onClick={() => setShowHistory(prev => !prev)}
            className="flex items-center gap-[10px] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors"
          >
            <span className="font-medium">
              История смен ({historyShifts.length})
            </span>

            <i
              className={`fa-solid fa-chevron-down text-[12px] transition-transform ${
                showHistory ? "rotate-180" : ""
              }`}
            />
          </button>

          <div
            className={`overflow-hidden transition-all duration-300 ${
              showHistory ? "max-h-[800px] mt-[16px]" : "max-h-0"
            }`}
          >
            <Table
              data={historyShifts}
              columns={columns}
              emptyMessage="История смен отсутствует"
            />
          </div>

        </div>
      )}

      {/* модалка */}
      <FormModal 
        key={editingShift?.shift ?? "new"}
        title={editingShift ? "Изменение смены" : "Добавление смены"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSave}
        initialValues={{
          dt: editingShift?.dt?.slice(0, 16) || ""
        }}
        fields={[
          {
            name: "dt",
            label: "Дата смены",
            type: "datetime",
            required: true,
          }
        ]}
      />

    </div>
  );
};

export default ShiftsPage;