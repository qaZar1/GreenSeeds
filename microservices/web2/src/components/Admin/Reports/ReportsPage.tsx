import React, { useEffect, useState } from "react";
import type { Report } from "../../../types/reports";
import type { Column } from "../../../types/table";
import { Table } from "../../utils/Table";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import { StatCard } from "../../utils/Card";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const ReportsPage: React.FC = () => {

  usePageHeader("Отчёты", "История выполнения заданий");

  const [reports, setReports] = useState<Report[]>([]);
  const [openShift, setOpenShift] = useState<number | null>(null);
  const [openHistoryShift, setOpenHistoryShift] = useState<number | null>(null);
  const [showHistory, setShowHistory] = useState(false);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadReports = async () => {
    setLoading(true);
    setError(false);

    try {
      const data = await api.getList("reports");
      setReports(data || []);
    } catch {
      setError(true);
      toast.error("Ошибка загрузки отчётов");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadReports();
  }, []);

  const today = new Date();
  today.setHours(0,0,0,0);

  const upcomingReports = reports.filter(
    r => new Date(r.dt || "") >= today
  );

  const historyReports = reports.filter(
    r => new Date(r.dt || "") < today
  );

  const columns: Column<Report>[] = [
    {
      header: "Задание",
      render: rec => (
        <div className="text-[14px] text-[var(--text-primary)]">
          {rec.number}
        </div>
      ),
    },
    {
      header: "Рецепт",
      render: rec => (
        <div className="text-[14px] text-[var(--text-primary)]">
          {rec.receipt}
        </div>
      ),
    },
    {
      header: "Выполнение",
      render: rec => (
        <div className="text-[14px] text-[var(--text-primary)]">
          {rec.turn}
        </div>
      ),
    },
    {
      header: "Дата",
      render: rec => (
        <div className="text-[14px] text-[var(--text-secondary)]">
          {rec.dt ? new Date(rec.dt).toLocaleString("ru-RU") : "-"}
        </div>
      ),
    },
    {
      header: "Статус",
      render: rec => (
        <span
          className={`px-[8px] py-[4px] rounded-[6px] text-[12px] font-medium ${
            rec.success
              ? "bg-[var(--status-success-bg)] text-[var(--status-success-text)]"
              : "bg-[var(--status-danger-bg)] text-[var(--status-danger-text)]"
          }`}
        >
          {rec.success ? "Успешно" : "Ошибка"}
        </span>
      ),
    },
    {
      header: "Действия",
      width: "80px",
      headerClassName: "text-right",
      className: "text-right",
      render: rec => (
        <button
          onClick={() => window.location.href = `/reports/${rec.id}`}
          className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:text-[var(--status-info-text)] hover:bg-[var(--status-info-bg)] transition-colors"
          title="Открыть"
        >
          <i className="fa-solid fa-eye text-[14px]" />
        </button>
      ),
    },
  ];

  const groupByShift = (list: Report[]) =>
    list.reduce<Record<number, Report[]>>((acc, item) => {
      if (!acc[item.shift]) acc[item.shift] = [];
      acc[item.shift].push(item);
      return acc;
    }, {});

  const grouped = groupByShift(upcomingReports);
  const groupedHistory = groupByShift(historyReports);

  const shifts = Object.keys(grouped).map(Number).sort((a,b) => b-a);
  const historyShifts = Object.keys(groupedHistory).map(Number).sort((a,b) => b-a);

  
  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadReports} />;

  return (
    <div className="space-y-[24px] w-full">

      {/* статистика */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-[12px]">
        <StatCard
          title="Всего отчётов"
          value={reports.length}
        />

        <StatCard
          title="Смен"
          value={shifts.length}
          color="var(--status-info-text)"
        />
      </div>

      {/* актуальные смены */}
      {shifts.map(shift => {

        const items = grouped[shift].sort((a,b) => {
          if (a.number !== b.number) return a.number - b.number;
          return a.turn - b.turn;
        });

        const isOpen = openShift === shift;

        return (
          <div
            key={shift}
            className="border border-[var(--border-color)] rounded-[12px] overflow-hidden"
          >

            <button
              onClick={() => setOpenShift(isOpen ? null : shift)}
              className="w-full flex items-center justify-between px-[16px] py-[12px] bg-[var(--bg-secondary)] hover:bg-[var(--bg-hover)] transition-colors"
            >
              <span className="font-medium text-[var(--text-primary)]">
                Смена №{shift}
              </span>

              <i
                className={`fa-solid fa-chevron-down text-[12px] transition-transform text-[var(--text-primary)]${
                  isOpen ? "rotate-180" : ""
                }`}
              />
            </button>

            {isOpen && (
              <div className="p-[12px]">
                <Table
                  data={items}
                  columns={columns}
                  emptyMessage="Отчёты отсутствуют"
                />
              </div>
            )}

          </div>
        );
      })}

      {/* история */}
      {historyShifts.length > 0 && (
        <div className="border-t border-[var(--border-color)] pt-[16px]">

          <button
            onClick={() => setShowHistory(prev => !prev)}
            className="flex items-center gap-[10px] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors"
          >
            <span className="font-medium">
              История отчётов ({historyReports.length})
            </span>

            <i
              className={`fa-solid fa-chevron-down text-[12px] transition-transform ${
                showHistory ? "rotate-180" : ""
              }`}
            />
          </button>

          <div
            className={`space-y-[16px] overflow-hidden transition-all duration-300 ${
              showHistory ? "max-h-[2000px] mt-[16px]" : "max-h-0"
            }`}
          >

            {historyShifts.map(shift => {

              const items = groupedHistory[shift].sort((a,b) => {
                if (a.number !== b.number) return a.number - b.number;
                return a.turn - b.turn;
              });

              const isOpen = openHistoryShift === shift;

              return (
                <div
                  key={shift}
                  className="border border-[var(--border-color)] rounded-[12px] overflow-hidden"
                >

                  <button
                    onClick={() =>
                      setOpenHistoryShift(isOpen ? null : shift)
                    }
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
                        emptyMessage="Отчёты отсутствуют"
                      />
                    </div>
                  )}

                </div>
              );
            })}

          </div>

        </div>
      )}

    </div>
  );
};

export default ReportsPage;