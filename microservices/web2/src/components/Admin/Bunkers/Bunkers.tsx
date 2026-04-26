import React, { useEffect, useState } from "react";
import type { Bunker } from "../../../types/bunker";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import FormModal from "../../utils/FormModal";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";


const BunkersPage: React.FC = () => {
  usePageHeader("Бункеры", "Настройка бункеров и расстояния подачи");

  const [bunkers, setBunkers] = useState<Bunker[]>([]);
  const [editingBunker, setEditingBunker] = useState<Bunker | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadBunkers = async () => {
    setLoading(true);
    setError(false);

    try {
      const data = await api.getList("bunkers");
      setBunkers(data || []);
    } catch {
      toast.error("Ошибка загрузки бункеров");
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadBunkers();
  }, []);

  const today = new Date();
  today.setHours(0, 0, 0, 0);


  const handleAddBunker = () => {
    setEditingBunker(null);
    setIsModalOpen(true);
  };

  const handleEdit = (bunker: Bunker) => {
    setEditingBunker(bunker);
    setIsModalOpen(true);
  };

  const handleSave = async (data: Record<string, any>) => {
    const loading = toast.loading(
      editingBunker ? "Изменение бункера" : "Добавление бункера"
    );

    try {
      if (editingBunker) {
        const updatedBunker = await api.update("bunkers", {
          id: editingBunker.bunker,
          distance: data.distance,
        });

        setBunkers(prev =>
          prev.map(b =>
            b.bunker === updatedBunker.data.bunker ? updatedBunker.data : b
          )
        );

        toast.success("Данные бункера обновлены", { id: loading });
      } else {
        const newBunker = await api.create("bunkers", {
          bunker: data.bunker,
          distance: data.distance
        });

        setBunkers(prev => [...prev, newBunker.data]);

        toast.success("Бункер создан", { id: loading });
      }

      setEditingBunker(null);
      setIsModalOpen(false);
    } catch {
      toast.error("Ошибка сохранения изменений", { id: loading });
    }
  };

  const handleDelete = async (bunker: Bunker) => {
    const loading = toast.loading("Удаление информации о бункере...");

    try {
      await api.delete("bunkers", bunker.bunker);

      setBunkers(prev => prev.filter(b => b.bunker !== bunker.bunker));

      toast.success("Бункер удален", { id: loading });
    } catch {
      toast.error("Ошибка удаления бункера", { id: loading });
    }
  };

  const columns: Column<Bunker>[] = [
    {
      header: "Бункер",
      render: bunker => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {bunker.bunker}
        </div>
      ),
    },
    {
      header: "Дистанция",
      render: bunker => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {bunker.distance}
        </div>
      ),
    },
    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: shift => (
        <div className="flex items-center justify-end gap-[8px]">

          <button
            onClick={() => handleEdit(shift)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)] transition-colors"
            title="Изменить данные о бункере"
          >
            <i className="fa-solid fa-pen-to-square text-[14px]" />
          </button>

          <button
            onClick={() => handleDelete(shift)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)] transition-colors"
            title="Удалить бункер"
          >
            <i className="fa-solid fa-trash text-[14px]" />
          </button>

        </div>
      ),
    },
  ];
  
  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadBunkers} />

  return (
    <div className="space-y-[24px] w-full">

      {/* кнопка добавления */}
      <div className="flex justify-end">
        <button
          onClick={handleAddBunker}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px] font-medium hover:bg-[var(--color-primary-hover)] transition-colors shadow-sm"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>
      </div>

      <Table
        data={bunkers}
        columns={columns}
        emptyMessage="Бункеры еще не заданы"
      />

      {/* модалка */}
      <FormModal 
        key={editingBunker?.bunker ?? "new"}
        title={editingBunker ? "Изменение бункера" : "Добавление бункера"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSave}
        initialValues={{
          bunker: editingBunker?.bunker || "",
          distance: editingBunker?.distance || ""
        }}
        fields={[
          {
            name: "bunker",
            label: "Номер бункера",
            type: "number",
            required: true,
            disabled: editingBunker ? true : false,
            min: 0,
          },
          {
            name: "distance",
            label: "Дистанция",
            type: "number",
            required: true,
            min: 0,
          }
        ]}
      />

    </div>
  );
};

export default BunkersPage;