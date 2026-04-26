import React, { useEffect, useMemo, useState } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";

import { Table } from "../../utils/Table";
import { StatCard } from "../../utils/Card";

import type { Column } from "../../../types/table";
import type { Seed } from "../../../types/seed";
import type { Receipt } from "../../../types/receipt";
import { useNavigate } from "react-router-dom";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const ReceiptPage: React.FC = () => {

  usePageHeader(
    "Рецепты",
    "Рецепты выращивания семян"
  );

  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [loading, setLoading] = useState(true);
  const [seeds, setSeeds] = useState<Seed[]>([]);
  const [error, setError] = useState(false);

  /* ---------------- загрузка ---------------- */

  const loadData = async () => {
    setLoading(true);
    setError(false);

    try {
      const [receiptsRes, seedsRes] = await Promise.allSettled([
        api.getList("receipts"),
        api.getList("seeds")
      ]);

      let hasError = false;

      if (receiptsRes.status === "fulfilled") {
        setReceipts(receiptsRes.value || []);
      } else {
        toast.error("Не удалось загрузить рецепты");
        hasError = true;
      }

      if (seedsRes.status === "fulfilled") {
        setSeeds(seedsRes.value || []);
      } else {
        toast.error("Не удалось загрузить семена");
        hasError = true;
      }

      if (hasError) {
        setError(true);
      }

    } catch {
      // сюда почти не зайдёт, но на всякий случай
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  /* ---------------- lookup ---------------- */

  const seedMap = useMemo(
    () => Object.fromEntries(seeds.map(s => [s.seed, s])),
    [seeds]
  );

  /* ---------------- actions ---------------- */
  const handleDelete = async (receipt: Receipt) => {

    const loading = toast.loading("Удаление рецепта...");

    try {
      await api.delete("receipts", receipt.receipt);

      setReceipts(prev =>
        prev.filter(r => r.receipt !== receipt.receipt)
      );

      toast.success("Рецепт удалён", { id: loading });
    } catch {
      toast.error("Ошибка удаления", { id: loading });
    }

  };

  /* ---------------- columns ---------------- */

  const columns: Column<Receipt>[] = [

    {
      header: "Семена",
      render: r => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {seedMap[r.seed]?.seed_ru ?? r.seed}
        </div>
      )
    },

    {
      header: "Описание",
      render: r => (
        <div className="text-[var(--text-primary)] text-[14px]">
          {r.description}
        </div>
      )
    },

    {
      header: "Обновлено",
      render: r => (
        <div className="text-[var(--text-secondary)] text-[14px]">
          {r.updated
            ? new Date(r.updated).toLocaleString("ru-RU")
            : "-"}
        </div>
      )
    },

    {
      header: "Действия",
      width: "120px",
      headerClassName: "text-right",
      className: "text-right",
      render: r => (
        <div className="flex justify-end gap-[8px]">

          <button
            onClick={() => navigate(`/settings/receipts/${r.receipt}/edit`)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)]"
          >
            <i className="fa-solid fa-pen-to-square text-[14px]" />
          </button>

          <button
            onClick={() => handleDelete(r)}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-danger-text)] hover:bg-[var(--status-danger-bg)]"
          >
            <i className="fa-solid fa-trash text-[14px]" />
          </button>

        </div>

      )
    }

  ];

  /* ---------------- render ---------------- */

	const navigate = useNavigate();

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadData} />;
    
  return (

    <div className="space-y-[24px] w-full">

      {/* add */}

      <div className="flex justify-end">

        <button
					onClick={() => navigate("/settings/receipts/create")}
          className="inline-flex items-center gap-[8px] px-[20px] py-[10px] bg-[var(--color-primary)] text-[var(--text-inverse)] rounded-[10px]"
        >
          <i className="fa-solid fa-plus text-[14px]" />
          Добавить
        </button>

      </div>

      {/* stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-[12px]">
        <StatCard
          title="Всего рецептов"
          value={receipts.length}
        />
        <StatCard
          title="Обновлены сегодня"
          value={
            receipts.filter(r => {
              if (!r.updated) return false;

              const today = new Date().toDateString();
              return new Date(r.updated).toDateString() === today;
            }).length
          }
          color="var(--status-info-text)"
        />

      </div>

      {/* table */}
      <Table
        data={receipts}
        columns={columns}
        emptyMessage="Рецепты ещё не созданы"
      />
    </div>

  );

};

export default ReceiptPage;