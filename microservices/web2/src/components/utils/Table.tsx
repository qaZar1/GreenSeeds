import { useState, useMemo } from "react"
import type { Column } from "../../types/table";

interface TableProps<T> {
  data: T[];
  columns: Column<T>[];
  emptyMessage?: string;
}

export function Table<T>({
  data,
  columns,
  emptyMessage = "Нет данных",
}: TableProps<T>) {
  return (
    <div className="overflow-x-auto rounded-[12px] border border-[var(--border-color)] bg-[var(--bg-card)] shadow-sm">
      <table className="w-full text-left">
        <colgroup>
          {columns.map((col, i) => (
            <col key={i} style={{ width: col.width }} />
          ))}
        </colgroup>
        
        {/* HEADER */}
        <thead className="bg-[var(--bg-hover)] border-b border-[var(--border-color)]">
          <tr>
            {columns.map((col, i) => (
              <th
                key={i}
                className={`px-[20px] py-[14px] text-[12px] font-semibold text-[var(--text-secondary)] uppercase tracking-wide ${col.headerClassName || ""}`}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>

        {/* BODY */}
        <tbody className="divide-y divide-[var(--border-light)]">
          {data.map((row, rowIndex) => (
            <tr
              key={rowIndex}
              className="hover:bg-[var(--bg-hover)] transition-colors"
            >
              {columns.map((col, colIndex) => (
                <td
                  key={colIndex}
                  className={`px-[20px] py-[16px] ${col.className || ""}`}
                >
                  {col.render(row)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>

      {data.length === 0 && (
        <div className="px-[20px] py-[40px] text-center text-[var(--text-secondary)]">
          <p className="text-[14px]">{emptyMessage}</p>
        </div>
      )}
    </div>
  );
}