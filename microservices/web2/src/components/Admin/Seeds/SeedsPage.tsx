import React, { useEffect, useState } from "react";
import type { Seed } from "../../../types/seed";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import FormModal from "../../utils/FormModal";
import { StatCard } from "../../utils/Card";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const SeedsPage: React.FC = () => {

  usePageHeader(
    "Семена",
    "Настройка параметров семян и емкости бункеров"
  );

  const [seeds, setSeeds] = useState<Seed[]>([]);
  const [editingSeed, setEditingSeed] = useState<Seed | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [showDeleted, setShowDeleted] = useState(false);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadSeeds = async () => {
    setLoading(true);
    setError(false);

    try {
      const data = await api.getList("seeds");
      setSeeds(data || []);
    } catch {
      setError(true);
      toast.error("Ошибка загрузки семян");
    } finally {
      setLoading(false);
    }
  };


  useEffect(() => {
    loadSeeds();
  }, []);

  /* -------- разделение активных и удаленных -------- */

  const activeSeeds = seeds.filter(s => !s.deleted_at);
  const deletedSeeds = seeds.filter(s => s.deleted_at);

  /* -------- actions -------- */

  const handleAdd = () => {
    setEditingSeed(null);
    setIsModalOpen(true);
  };

  const handleEdit = (seed: Seed) => {
    setEditingSeed(seed);
    setIsModalOpen(true);
  };

  const handleSave = async (data: Record<string, any>) => {

    const loading = toast.loading(
      editingSeed ? "Изменение семян" : "Добавление семян"
    );

    try {

      if (editingSeed) {

        const updated = await api.update("seeds", {
          id: editingSeed.seed,
          ...data
        });

        setSeeds(prev =>
          prev.map(s =>
            s.seed === updated.data.seed ? updated.data : s
          )
        );

        toast.success("Семена обновлены", { id: loading });

      } else {
        const created = await api.create("seeds", data);
        setSeeds(prev => [...prev, created.data]);
        toast.success("Семена добавлены", { id: loading });
      }

      setEditingSeed(null);
      setIsModalOpen(false);

    } catch {

      toast.error("Ошибка сохранения", { id: loading });

    }

  };

  const handleDelete = async (seed: Seed) => {

    const loading = toast.loading("Удаление семян...");

    try {
      await api.delete("seeds", seed.seed);

      setSeeds(prev =>
        prev.map(s =>
          s.seed === seed.seed
            ? { ...s, deleted_at: new Date().toISOString() }
            : s
        )
      );

      toast.success("Семена удалены", { id: loading });

    } catch {

      toast.error("Ошибка удаления", { id: loading });

    }

  };

  const handleRestore = async (seed: Seed) => {
    const loading = toast.loading("Восстановление...");

    try {
      const restored = await api.update("seeds", {
        id: seed.seed,
        ...seed,
        deleted_at: null,
      });


      setSeeds(prev =>
        prev.map(s =>
          s.seed === restored.data.seed ? restored.data : s
        )
      );

      toast.success("Семена восстановлены", { id: loading });
    } catch {
      toast.error("Ошибка восстановления", { id: loading });
    }
  };

  /* -------- columns -------- */

  const columns: Column<Seed>[] = [

    {
      header: "Семена",
      render: seed => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {seed.seed_ru}
        </div>
      )
    },

    {
      header: "Емкость бункера",
      render: seed => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {seed.tank_capacity}
        </div>
      )
    },

    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: seed => (

        <div className="flex items-center justify-end gap-[8px]">

          <button
            onClick={() => handleEdit(seed)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)] transition-colors"
          >
            <i className="fa-solid fa-pen-to-square text-[14px]" />
          </button>

          <button
            onClick={() => handleDelete(seed)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)] transition-colors"
          >
            <i className="fa-solid fa-trash text-[14px]" />
          </button>

        </div>

      )
    }

  ];

  const deletedColumns: Column<Seed>[] = [
    {
      header: "Семена",
      render: seed => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {seed.seed_ru}
        </div>
      )
    },

    {
      header: "Емкость",
      render: seed => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {seed.tank_capacity}
        </div>
      )
    },

    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: seed => (

        <div className="flex justify-end">
          <button
            onClick={() => handleRestore(seed)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-success-text)] hover:bg-[var(--status-success-bg)]"
          >
            <i className="fa-solid fa-rotate-left text-[14px]" />
          </button>
        </div>
      )
    }
  ];

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadSeeds} />;

  return (
    <div className="space-y-[24px] w-full">

      {/* add button */}
      <div className="flex justify-end">
        <button
          onClick={handleAdd}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px]"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>
      </div>

      {/* active */}
      <Table
        data={activeSeeds}
        columns={columns}
        emptyMessage="Семена еще не заданы"
      />

      {/* deleted */}
      {deletedSeeds.length > 0 && (
        <div className="border-t border-[var(--border-color)] pt-[16px]">
          <button
            onClick={() => setShowDeleted(prev => !prev)}
            className="flex items-center gap-[10px] text-[var(--text-secondary)]"
          >

            <span className="font-medium">
              Удаленные ({deletedSeeds.length})
            </span>
            <i
              className={`fa-solid fa-chevron-down text-[12px] transition-transform ${
                showDeleted ? "rotate-180" : ""
              }`}
            />
          </button>
          <div
            className={`overflow-hidden transition-all duration-300 ${
              showDeleted ? "max-h-[800px] mt-[16px]" : "max-h-0"
            }`}
          >
            <Table
              data={deletedSeeds}
              columns={deletedColumns}
              emptyMessage="Удаленных записей нет"
            />
          </div>
        </div>
      )}

      {/* modal */}

      <FormModal
        key={editingSeed?.seed ?? "new"}
        title={editingSeed ? "Изменение семян" : "Добавление семян"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSave}
        initialValues={{
          seed: editingSeed?.seed || "",
          seed_ru: editingSeed?.seed_ru || "",
          min_density: editingSeed?.min_density || "",
          max_density: editingSeed?.max_density || "",
          tank_capacity: editingSeed?.tank_capacity || ""
        }}
        fields={[
          {
            name: "seed",
            label: "Название (англ)",
            type: "text",
            required: true,
            disabled: !!editingSeed
          },
          {
            name: "seed_ru",
            label: "Название (рус)",
            type: "text",
            required: true
          },
          {
            name: "min_density",
            label: "Мин. плотность",
            type: "number",
            required: true
          },
          {
            name: "max_density",
            label: "Макс. плотность",
            type: "number",
            required: true
          },
          {
            name: "tank_capacity",
            label: "Емкость бункера",
            type: "number",
            required: true
          }
        ]}
      />
    </div>

  );

};

export default SeedsPage;