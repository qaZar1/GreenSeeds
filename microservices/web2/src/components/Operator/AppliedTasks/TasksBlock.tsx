import React from "react";
import { useNavigate } from "react-router-dom";

type Task = {
  id: number;
  shift: number;
  number: number;
  amount: number;
  seed_ru?: string;
};

const TasksBlock: React.FC<{ tasks: Task[] }> = ({ tasks }) => {
  const navigate = useNavigate();

  if (!tasks.length) {
    return (
      <div className="text-center py-[40px] text-[var(--text-secondary)]">
        Заданий пока нет
      </div>
    );
  }

  const grouped = tasks.reduce<Record<string, Task[]>>((acc, t) => {
    const key = t.shift ? `Смена №${t.shift}` : "Без смены";
    if (!acc[key]) acc[key] = [];
    acc[key].push(t);
    return acc;
  }, {});

  return (
    <div className="space-y-[20px]">

      {Object.entries(grouped).map(([shift, group]) => (
        <div
          key={shift}
          className="border border-[var(--border-color)] rounded-[12px] overflow-hidden text-[var(--text-primary)]"
        >

          {/* header */}
          <div className="px-[16px] py-[12px] bg-[var(--bg-secondary)] border-b border-[var(--border-color)] font-medium">
            {shift}
          </div>

          {/* table */}
          <div className="divide-y divide-[var(--border-color)]">

            {/* header row */}
            <div className="grid grid-cols-3 px-[16px] py-[10px] text-[13px] text-[var(--text-secondary)]">
              <div>Задание</div>
              <div>Кол-во</div>
              <div>Культура</div>
            </div>

            {/* rows */}
            {group.map((t) => (
              <div
                key={`${t.id}-${t.shift}`}
                onClick={() => navigate(`/tasks/${t.id}`)}
                className="grid grid-cols-3 px-[16px] py-[12px] text-[14px] cursor-pointer hover:bg-[var(--bg-hover)] transition"
              >
                <div className="text-[var(--text-primary)]">
                  {t.number}
                </div>

                <div className="text-[var(--text-primary)]">
                  {t.amount}
                </div>

                <div className="text-[var(--text-primary)]">
                  {t.seed_ru || "—"}
                </div>
              </div>
            ))}

          </div>
        </div>
      ))}

    </div>
  );
};

export default TasksBlock;