import React from "react";

export const StatCard = ({
  title,
  value,
  color,
}: {
  title: string;
  value: number | string;
  color?: string;
}) => (
  <div className="bg-[var(--bg-card)] rounded-[12px] p-[16px] border border-[var(--border-color)] text-center">
    <div
      className="text-[28px] font-bold"
      style={{ color: color || "var(--text-primary)" }}
    >
      {value}
    </div>
    <div className="text-[12px] text-[var(--text-secondary)]">{title}</div>
  </div>
);