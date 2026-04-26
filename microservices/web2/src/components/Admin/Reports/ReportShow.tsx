import React, { useEffect, useState } from "react";
import { useParams, useNavigate, data } from "react-router-dom";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import type { Report } from "../../../types/reports";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const Field: React.FC<{ label: string; value?: any }> = ({ label, value }) => (
  <div className="flex flex-col gap-[4px]">
    <span className="text-[12px] text-[var(--text-secondary)]">
      {label}
    </span>
    <span className="text-[14px] text-[var(--text-primary)]">
      {value ?? "-"}
    </span>
  </div>
);

const ReportShow: React.FC = () => {

  const { id } = useParams();
  const navigate = useNavigate();

  usePageHeader("Отчёт", "Детали выполнения задания");

  const [report, setReport] = useState<Report | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadReport = async () => {
    if (!id) return;

    setLoading(true);
    setError(false);

    try {
      const data = await api.getOne("reports", id);

      setReport(data?.data || null);
    } catch {
      setError(true);
      toast.error("Ошибка загрузки отчёта");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadReport();
  }, [id]);
  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadReport} />;

  if (!report) return "";

  return (
    <div className="space-y-[24px] w-full">

      {/* кнопка назад */}
      <div>
        <button
          onClick={() => navigate(-1)}
          className="px-[16px] py-[8px] rounded-[8px] border border-[var(--border-color)] hover:bg-[var(--bg-hover)] transition-colors text-[var(--text-primary)]"
        >
          ← Назад
        </button>
      </div>

      <div className="border border-[var(--border-color)] rounded-[12px] p-[20px] space-y-[20px]">

        <h2 className="text-[18px] font-semibold text-[var(--text-primary)]">
          Детали отчёта
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-[24px]">

          {/* левая колонка */}
          <div className="space-y-[16px]">

            <Field label="Смена" value={report.shift} />

            <Field label="Номер задания" value={report.number} />

            <Field label="Рецепт" value={report.receipt} />

            <Field label="Номер выполнения" value={report.turn} />

            <Field
              label="Дата"
              value={
                report.dt
                  ? new Date(report.dt).toLocaleString("ru-RU")
                  : "-"
              }
            />

          </div>

          {/* правая колонка */}
          <div className="space-y-[16px]">

            <div className="flex flex-col gap-[4px]">
              <span className="text-[12px] text-[var(--text-secondary)]">
                Успешно
              </span>

              <span
                className={`px-[8px] py-[4px] rounded-[6px] text-[12px] font-medium w-fit ${
                  report.success
                    ? "bg-[var(--status-success-bg)] text-[var(--status-success-text)]"
                    : "bg-[var(--status-danger-bg)] text-[var(--status-danger-text)]"
                }`}
              >
                {report.success ? "Да" : "Нет"}
              </span>
            </div>

            <Field label="Ошибка" value={report.error} />

            <Field label="Решение" value={report.solution} />

            <Field label="Маркировка" value={report.mark} />

            <Field
              label="Ответственный"
              value={report.responsible}
            />

          </div>

        </div>

      </div>

    </div>
  );
};

export default ReportShow;