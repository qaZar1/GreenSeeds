import React, { useState, useEffect, useCallback, useRef } from "react";
import { usePageHeader } from "../../../context/HeaderContext";
import { api } from "../../../api/apiProvider";
import { Table } from "../../utils/Table";
import type { Column } from "../../../types/table";
import SproutLoader from "../../utils/Loader/SproutLoader";
import toast from "react-hot-toast";
import ErrorState from "../../pages/ErrorState";
import ResponsiveTable from "../../utils/ResponsiveTable";

const LIMIT = 50;

const LogsPage: React.FC = () => {
  usePageHeader("Логи системы", "Журнал событий приложения");

  const scrollRef = useRef<HTMLDivElement | null>(null);
  const loadMoreRef = useRef<HTMLDivElement | null>(null);

  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [more, setMore] = useState(true);

  const [search, setSearch] = useState("");
  const [level, setLevel] = useState("ALL");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");

  const [offset, setOffset] = useState(0);
  const [query, setQuery] = useState("");

  const [openId, setOpenId] = useState<number | null>(null);

  const [error, setError] = useState(false);

  const toggleRow = (id: number) => {
    setOpenId(prev => (prev === id ? null : id));
  };

  useEffect(() => {
    const t = setTimeout(() => setQuery(search), 300);
    return () => clearTimeout(t);
  }, [search]);

  const fetchLogs = useCallback(async () => {
    if (loading) return;

    setLoading(true);
    setError(false);

    try {
      if (offset === 0) {
        setData([]);
        setMore(true);
      }

      const params = {
        limit: LIMIT,
        offset,
        ...(query && { search: query }),
        ...(level !== "ALL" && { level }),
        ...(dateFrom && { date_from: dateFrom }),
        ...(dateTo && { date_to: dateTo }),
      };

      const newDataRaw = (await api.getList("logs", params)) || [];

      if (newDataRaw.length === 0) {
        setMore(false);
        return;
      }

      const newData = newDataRaw.map((item: any, index: number) => ({
        id: offset + index,
        ...item,
      }));

      setData(prev => (offset === 0 ? newData : [...prev, ...newData]));
      setMore(newDataRaw.length === LIMIT);

    } catch (e) {
      console.error(e);
      toast.error("Не удалось загрузить логи");
      setError(true);
      setMore(false);
    } finally {
      setLoading(false);
    }
  }, [offset, query, level, dateFrom, dateTo, loading]);

  useEffect(() => {
    fetchLogs();
  }, [offset, query, level, dateFrom, dateTo]);

  // infinite scroll
  useEffect(() => {
    const el = loadMoreRef.current;
    if (!el) return;

    const observer = new IntersectionObserver(
      entries => {
        const entry = entries[0];

        if (
          entry.isIntersecting &&
          more &&
          !loading &&
          data.length >= LIMIT
        ) {
          setOffset(prev => prev + LIMIT);
        }
      },
      {
        root: null,
        rootMargin: "200px",
      }
    );

    observer.observe(el);

    return () => observer.disconnect();
  }, [more, loading]);

  const columns: Column<any>[] = [
    {
      header: "Дата",
      width: "180px",
      render: row => (
        <span className="text-[var(--text-secondary)]">
          {new Date(row.dt).toLocaleString("ru-RU")}
        </span>
      ),
    },
    {
      header: "Уровень",
      width: "120px",
      render: row => {
        const colors: Record<string, string> = {
          INFO: "bg-[var(--status-info-bg)] text-[var(--status-info-text)]",
          WARN: "bg-[var(--status-warning-bg)] text-[var(--status-warning-text)]",
          ERROR: "bg-[var(--status-danger-bg)] text-[var(--status-danger-text)]",
        };

        return (
          <span
            className={`px-[8px] py-[4px] rounded-[6px] text-[12px] font-medium ${colors[row.lvl] || ""}`}
          >
            {row.lvl}
          </span>
        );
      },
    },
    {
      header: "Сообщение",
      render: row => (
        <div
          onClick={() => toggleRow(row.id)}
          className="cursor-pointer text-[var(--text-primary)] text-[14px]"
        >
          {row.msg}

          {openId === row.id && (
            <div className="mt-[10px] p-[12px] rounded-[8px] bg-[var(--bg-secondary)] text-[13px] space-y-[4px] border border-[var(--border-color)]">
              <div><b>Дата:</b> {row.dt}</div>
              <div><b>ID запроса:</b> {row.request_id || "-"}</div>
              <div><b>Функция:</b> {row.caller || "-"}</div>
              <div><b>Пользователь:</b> {row.username || "-"}</div>
              <div><b>Сообщение:</b> {row.msg || "-"}</div>
            </div>
          )}
        </div>
      ),
    },
  ];

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={fetchLogs}/>

  return (
    <div className="flex flex-col w-full space-y-[20px]">
      {/* фильтры */}
      <div className="grid gap-[12px] grid-cols-1 sm:grid-cols-2 md:grid-cols-4">
        <input
          value={search}
          onChange={e => {
            setOffset(0);
            setSearch(e.target.value);
          }}
          placeholder="Поиск"
          className="px-[14px] py-[10px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)]"
        />

        <select
          value={level}
          onChange={e => {
            setOffset(0);
            setLevel(e.target.value);
          }}
          className="px-[14px] py-[10px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)]"
        >
          <option value="ALL">Все</option>
          <option value="INFO">INFO</option>
          <option value="WARN">WARN</option>
          <option value="ERROR">ERROR</option>
        </select>

        <input
          type="date"
          value={dateFrom}
          onChange={e => {
            setOffset(0);
            setDateFrom(e.target.value);
          }}
          className="date px-[14px] py-[10px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)]"
        />

        <input
          type="date"
          value={dateTo}
          onChange={e => {
            setOffset(0);
            setDateTo(e.target.value);
          }}
          className="date px-[14px] py-[10px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)]"
        />
      </div>

      {/* таблица */}
<div
  ref={scrollRef}
  className="overflow-y-auto flex-1 min-h-0 bg-[var(--bg-page)]"
>
  <ResponsiveTable
    data={data}
    table={
      <Table
        data={data}
        columns={columns}
        emptyMessage={loading ? "Загрузка..." : "Логи отсутствуют"}
      />
    }
    emptyMessage={loading ? "Загрузка..." : "Логи отсутствуют"}
    renderCard={(row) => {
      const colors: Record<string, string> = {
        INFO: "bg-[var(--status-info-bg)] text-[var(--status-info-text)]",
        WARN: "bg-[var(--status-warning-bg)] text-[var(--status-warning-text)]",
        ERROR: "bg-[var(--status-danger-bg)] text-[var(--status-danger-text)]",
      };

      const isOpen = openId === row.id;

      return (
        <>
          {/* content */}
          <div className="space-y-[10px] text-[14px]">

            <div className="text-[var(--text-primary)]">
              <span className="text-[var(--text-secondary)]">
                Дата:
              </span>{" "}
              {new Date(row.dt).toLocaleString("ru-RU")}
            </div>

            <div className="text-[var(--text-primary)]">
              <span className="text-[var(--text-secondary)]">
                Уровень:
              </span>{" "}

              <span
                className={`inline-flex px-[8px] py-[4px] rounded-[6px] text-[12px] font-medium ${colors[row.lvl] || ""}`}
              >
                {row.lvl}
              </span>
            </div>

            <div className="text-[var(--text-primary)] break-words">
              <span className="text-[var(--text-secondary)]">
                Сообщение:
              </span>{" "}
              {row.msg}
            </div>

            {isOpen && (
              <div className="pt-[4px] space-y-[6px] text-[13px]">
                <div>
                  <span className="text-[var(--text-secondary)]">
                    ID запроса:
                  </span>{" "}
                  {row.request_id || "-"}
                </div>

                <div>
                  <span className="text-[var(--text-secondary)]">
                    Функция:
                  </span>{" "}
                  {row.caller || "-"}
                </div>

                <div>
                  <span className="text-[var(--text-secondary)]">
                    Пользователь:
                  </span>{" "}
                  {row.username || "-"}
                </div>
              </div>
            )}

          </div>

          {/* actions */}
          <button
            onClick={() => toggleRow(row.id)}
            className="
              w-full
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
            <i
              className={`fa-solid ${
                isOpen ? "fa-chevron-up" : "fa-chevron-down"
              }`}
            />

            {isOpen ? "Скрыть" : "Подробнее"}
          </button>
        </>
      );
    }}
  />

  {loading && offset > 0 && (
    <div className="text-center py-[16px] text-[var(--text-secondary)] text-[14px]">
      Загрузка...
    </div>
  )}

  <div ref={loadMoreRef} style={{ height: 1 }} />
</div>
    </div>
  );
};

export default LogsPage;