import React from 'react';

interface ResponsiveTableProps<T> {
  data: T[];

  table: React.ReactNode;

  renderCard: (item: T) => React.ReactNode;

  emptyMessage?: string;
}

const ResponsiveTable = <T,>({
  data,
  table,
  renderCard,
  emptyMessage = 'Нет данных',
}: ResponsiveTableProps<T>) => {
  if (!data.length) {
    return (
      <div
        className="
          bg-[var(--bg-card)]
          border border-[var(--border-color)]
          rounded-[12px]
          p-[20px]
          text-center
          text-[14px]
          text-[var(--text-secondary)]
        "
      >
        {emptyMessage}
      </div>
    );
  }

  return (
    <>
      {/* MOBILE CARDS */}
      <div className="lg:hidden space-y-[12px]">
        {data.map((item, index) => (
          <div
            key={index}
            className="
              bg-[var(--bg-card)]
              border border-[var(--border-color)]
              rounded-[12px]
              p-[14px]
              space-y-[14px]
              overflow-hidden
            "
          >
            {renderCard(item)}
          </div>
        ))}
      </div>

      {/* DESKTOP TABLE */}
      <div className="hidden lg:block">
        {table}
      </div>
    </>
  );
};

export default ResponsiveTable;