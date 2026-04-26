import React, { useEffect, useState } from "react";
import type { Assignment } from "../../../types/assignments";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import { StatCard } from "../../utils/Card";
import FormModal from "../../utils/FormModal";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const AssignmentsPage: React.FC = () => {
	usePageHeader("Сменные задания", "Управление производственными заданиями");

  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [shiftsList, setShiftsList] = useState<any[]>([]);
  const [receiptsList, setReceiptsList] = useState<any[]>([]);

  const [openShift, setOpenShift] = useState<number | null>(null);

  const [editingAssignment, setEditingAssignment] = useState<Assignment | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadData = async () => {
    setLoading(true);
    setError(false);

    try {
      const [assignmentsRes, shiftsRes, receiptsRes] = await Promise.all([
        api.getList("assignments"),
        api.getList("shifts"),
        api.getList("receipts"),
      ]);

      setAssignments(assignmentsRes);
      setShiftsList(shiftsRes);
      setReceiptsList(receiptsRes);

    } catch (e: any) {
      console.error(e);

      toast.error("Не удалось загрузить данные");
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAddAssignment = () => {
    setEditingAssignment(null);
    setIsModalOpen(true);
  };

  const handleEdit = (rec: Assignment) => {
    setEditingAssignment(rec);
    setIsModalOpen(true);
  };

  const handleDelete = async (rec: Assignment) => {
    const loading = toast.loading("Удаление задания");

    try {
      await api.delete("assignments", rec.id);

      setAssignments(prev => prev.filter(a => a.id !== rec.id));

      toast.success("Задание удалено", { id: loading });
    } catch {
      toast.error("Ошибка удаления", { id: loading });
    }
  };

  const handleSave = async (data: Record<string, any>) => {
    const loading = toast.loading(
      editingAssignment ? "Изменение задания" : "Добавление задания"
    );

    try {
      if (editingAssignment) {
				console.log(editingAssignment)
        const updated = await api.update("assignments", {
					id: editingAssignment.id,
          ...data,
        });

        console.log("UPDATED", updated);

        setAssignments(prev =>
          prev.map(a => (a.id === updated.data.id ? updated.data : a))
        );

        toast.success("Задание обновлено", { id: loading });

      } else {
        const created = await api.create("assignments", data);
        console.log("CREATED", created);
        setAssignments(prev => [...prev, created.data]);
        toast.success("Задание создано", { id: loading });
      }

      setEditingAssignment(null);
      setIsModalOpen(false);

    } catch {
      toast.error("Ошибка сохранения", { id: loading });
    }
  };

  const columns: Column<Assignment>[] = [
    {
      header: "№",
      render: rec => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {rec.number}
        </div>
      ),
    },
    {
      header: "Рецепт",
      render: rec => {
        const receipt = receiptsList.find(r => r.id === rec.receipt);

        return (
          <div className="text-[var(--text-primary)] text-[14px]">
            {receipt?.name ?? `Рецепт ${rec.receipt}`}
          </div>
        );
      },
    },
    {
      header: "Количество",
      render: rec => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {rec.amount}
        </div>
      ),
    },
    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: rec => (
        <div className="flex items-center justify-end gap-[8px]">

          <button
            onClick={() => handleEdit(rec)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)] transition-colors"
            title="Редактировать"
          >
            <i className="fa-solid fa-pen text-[14px]" />
          </button>

          <button
            onClick={() => handleDelete(rec)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)] transition-colors"
            title="Удалить"
          >
            <i className="fa-solid fa-trash text-[14px]" />
          </button>

        </div>
      ),
    },
  ];

  const grouped = assignments.reduce<Record<number, Assignment[]>>((acc, item) => {
    if (!acc[item.shift]) acc[item.shift] = [];
    acc[item.shift].push(item);
    return acc;
  }, {});

  const shifts = Object.keys(grouped).sort((a, b) => Number(a) - Number(b));

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadData} />

  return (
    <div className="space-y-[24px] w-full">

      {/* кнопка добавления */}
      <div className="flex justify-end">
        <button
          onClick={handleAddAssignment}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px] font-medium hover:bg-[var(--color-primary-hover)] transition-colors shadow-sm"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>
      </div>

      {/* статистика */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-[12px]">
        <StatCard title="Всего заданий" value={assignments.length} />
        <StatCard
          title="Смен"
          value={shifts.length}
          color="var(--status-success-text)"
        />
      </div>

      {/* список смен */}
      {shifts.map(shift => {
        const items = grouped[Number(shift)] || [];
        const isOpen = openShift === Number(shift);

        return (
          <div
            key={shift}
            className="border border-[var(--border-color)] rounded-[12px] overflow-hidden"
          >

            <button
              onClick={() => setOpenShift(isOpen ? null : Number(shift))}
              className="w-full flex items-center justify-between px-[16px] py-[12px] bg-[var(--bg-secondary)] hover:bg-[var(--bg-hover)] transition-colors"
            >
              <span className="font-medium text-[var(--text-primary)]">
                Смена №{shift}
              </span>

              <i
                className={`fa-solid fa-chevron-down text-[12px] transition-transform text-[var(--text-primary)] ${
                  isOpen ? "rotate-180" : ""
                }`}
              />
            </button>

            {isOpen && (
              <div className="p-[12px]">
                <Table
                  data={items}
                  columns={columns}
                  emptyMessage="Задания отсутствуют"
                />
              </div>
            )}

          </div>
        );
      })}

      {/* модалка */}
      <FormModal
        key={editingAssignment?.id ?? "new"}
        title={editingAssignment ? "Изменение задания" : "Добавление задания"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSave}
        initialValues={{
          shift: editingAssignment?.shift,
          number: editingAssignment?.number,
          receipt: editingAssignment?.receipt,
          amount: editingAssignment?.amount,
        }}
        fields={[
          {
            name: "shift",
            label: "Смена",
            type: "select",
            required: true,
            disabled: !!editingAssignment,
            options: [
							{ label: "Выберите смену", value: "" },
							...shiftsList
                .filter(s => {
                  const today = new Date();
                  today.setHours(0, 0, 0, 0);

                  return new Date(s.dt) >= today;
                })
                .map(s => ({
                  label: `Смена ${new Date(s.dt).toLocaleString("ru-RU")}`,
                  value: s.shift,
                }))
						],
          },
          {
            name: "number",
            label: "Номер задания",
            type: "number",
            required: true,
            disabled: !!editingAssignment,
          },
          {
            name: "receipt",
            label: "Рецепт",
            type: "select",
            required: true,
            disabled: !!editingAssignment,
            options: [
							{ label: "Выберите рецепт", value: "" },
							...receiptsList.map(r => ({
								label: r.seed_ru,
								value: r.id,
							}))
						],
          },
          {
            name: "amount",
            label: "Количество",
            type: "number",
            required: true,
          },
        ]}
      />

    </div>
  );
};

export default AssignmentsPage;