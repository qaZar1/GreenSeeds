import React, { useEffect, useState, useMemo } from "react";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import FormModal from "../../utils/FormModal";
import { StatCard } from "../../utils/Card";

import type { Placement } from "../../../types/placement";
import type { Seed } from "../../../types/seed";
import type { Bunker } from "../../../types/bunker";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const PlacementPage: React.FC = () => {

  usePageHeader(
    "Размещение семян",
    "Количество семян в бункерах"
  );

  const [placements, setPlacements] = useState<Placement[]>([]);
  const [editingPlacement, setEditingPlacement] = useState<Placement | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const [seeds, setSeeds] = useState<Seed[]>([]);
  const [bunkers, setBunkers] = useState<Bunker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  /* ---------------- загрузка ---------------- */
  const loadData = async () => {
    setLoading(true);
    setError(false);

    try {
      const [placementsRes, seedsRes, bunkersRes] = await Promise.allSettled([
        api.getList("placements"),
        api.getList("seeds"),
        api.getList("bunkers")
      ]);

      const allFulfilled =
        placementsRes.status === "fulfilled" &&
        seedsRes.status === "fulfilled" &&
        bunkersRes.status === "fulfilled";

      if (!allFulfilled) {
        toast.error("Не удалось загрузить данные");
        setError(true);
        return;
      }

      // если всё ок — записываем
      setPlacements(placementsRes.value || []);
      setSeeds(seedsRes.value || []);
      setBunkers(bunkersRes.value || []);

    } catch {
      toast.error("Не удалось загрузить данные");
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  /* ---------------- быстрые lookup структуры ---------------- */
  const sortedPlacements = useMemo(() => {
    return [...placements].sort((a, b) => a.bunker - b.bunker);
  }, [placements]);

  const seedMap = useMemo(
    () => Object.fromEntries(seeds.map(s => [s.seed, s])),
    [seeds]
  );

  const usedBunkers = useMemo(
    () => new Set(placements.map(p => p.bunker)),
    [placements]
  );

  /* ---------------- actions ---------------- */
	const editingBunker = editingPlacement?.bunker;

  const handleAdd = () => {
    setEditingPlacement(null);
    setIsModalOpen(true);
  };

  const handleEdit = (placement: Placement) => {
    setEditingPlacement(placement);
    setIsModalOpen(true);
  };

  const handleSave = async (data: Record<string, any>) => {
    const seed = seedMap[data.seed];

    if (seed && data.amount > seed.tank_capacity) {
      toast.error(`Максимум для этих семян: ${seed.tank_capacity}`);
      return;
    }

    const loading = toast.loading(
      editingPlacement ? "Изменение размещения" : "Добавление размещения"
    );

    try {
      if (editingPlacement) {
        const updated = await api.update("placements", {
          id: data.bunker,
          amount: data.amount,
          seed: data.seed
        });

        setPlacements(prev =>
          prev.map(p =>
            p.bunker === updated.data.bunker && p.seed === updated.data.seed
              ? updated.data
              : p
          )
        );
        toast.success("Размещение обновлено", { id: loading });
      } else {
        const created = await api.create("placements", data);
        setPlacements(prev => [...prev, created.data]);
        toast.success("Размещение добавлено", { id: loading });
      }

      setEditingPlacement(null);
      setIsModalOpen(false);
    } catch {
      toast.error("Ошибка сохранения", { id: loading });
    }

  };

  const handleDelete = async (placement: Placement) => {
    const loading = toast.loading("Удаление размещения...");

    try {
      await api.delete("placements", placement.bunker);

      setPlacements(prev =>
        prev.filter(p =>
          !(p.bunker === placement.bunker && p.seed === placement.seed)
        )
      );

      toast.success("Размещение удалено", { id: loading });

    } catch {
      toast.error("Ошибка удаления", { id: loading });
    }

  };

  /* ---------------- columns ---------------- */

  const columns: Column<Placement>[] = [

    {
      header: "Бункер",
      render: p => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {p.bunker}
        </div>
      )
    },

    {
      header: "Семена",
      render: p => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {seedMap[p.seed]?.seed_ru ?? p.seed}
        </div>
      )
    },

    {
      header: "Количество (шт. лотков)",
      render: p => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {p.amount}
        </div>
      )
    },

    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: p => (

        <div className="flex items-center justify-end gap-[8px]">

          <button
            onClick={() => handleEdit(p)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)] transition-colors"
          >
            <i className="fa-solid fa-pen-to-square text-[14px]" />
          </button>

          <button
            onClick={() => handleDelete(p)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)] transition-colors"
          >
            <i className="fa-solid fa-trash text-[14px]" />
          </button>

        </div>

      )
    }

  ];

  /* ---------------- render ---------------- */

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadData} />;

  return (

    <div className="space-y-[24px] w-full">

      {/* add */}

      <div className="flex justify-end">

        <button
          onClick={handleAdd}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px]"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>

      </div>

      {/* stats */}

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-[12px]">

        <StatCard
          title="Всего размещений"
          value={placements.length}
        />

        <StatCard
          title="Пустых бункеров"
          value={placements.filter(p => p.amount === 0).length}
          color="var(--status-danger-text)"
        />

        <StatCard
          title="Заканчиваются семена"
          value={
            placements.filter(p => {

              const seed = seedMap[p.seed];
              if (!seed) return false;

              const percent = p.amount / seed.tank_capacity;

              return percent > 0 && percent < 0.1;

            }).length
          }
          color="var(--status-warning-text)"
        />

      </div>

      {/* table */}

      <Table
        data={sortedPlacements}
        columns={columns}
        emptyMessage="Размещения еще не заданы"
      />

      {/* modal */}

      <FormModal
        key={editingPlacement ? `${editingPlacement.bunker}-${editingPlacement.seed}` : "new"}
        title={editingPlacement ? "Изменение размещения" : "Добавление размещения"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSave}
        initialValues={{
          bunker: editingPlacement?.bunker || "",
          seed: editingPlacement?.seed || "",
          amount: editingPlacement?.amount || 0
        }}
        fields={[
          {
						name: "bunker",
						label: "Бункер",
						type: "select",
						required: true,
						disabled: !!editingPlacement,
						options: [
							{ label: "Выберите бункер", value: "" },
							...bunkers
								.filter(b => b.bunker !== null)
								.filter(b => !usedBunkers.has(b.bunker) || b.bunker === editingBunker)
								.sort((a, b) => a.bunker - b.bunker)
								.map(b => ({
									label: `Бункер ${b.bunker}`,
									value: b.bunker
								}))
						]
					},
          {
            name: "seed",
            label: "Семена",
            type: "select",
            required: true,
            options: [
              { label: "Выберите семена", value: "" },
              ...seeds.map(s => ({
                label: `${s.seed_ru}`,
                value: s.seed
              }))
            ]
          },
          {
            name: "amount",
            label: "Количество",
            type: "number",
            required: true,
            min: 0,
						max: (values) => seedMap[values.seed]?.tank_capacity
          }
        ]}
      />

    </div>

  );

};

export default PlacementPage;