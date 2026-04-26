import React, { useState } from "react";
import { api } from "../../../api/apiProvider";
import toast from "react-hot-toast";
import { useAuth } from "../../../context/AuthContext";

type Props = {
  task: {
    id: number;
    shift: number;
    dt: string;
  };
  onTaken: () => void;
};

const TaskCard: React.FC<Props> = ({ task, onTaken }) => {
  const [loading, setLoading] = useState(false);

  const date = new Date(task.dt);
  const isStarted = date <= new Date();

  const { user_id } = useAuth();

  const handleTake = async () => {
    const t = toast.loading("Берём задание...");
    setLoading(true);

    try {
      await api.update("choice", {
        shift: task.shift,
        user_id: user_id,
        dt: task.dt,
      });

      toast.success("Задание взято", { id: t });
      onTaken();
    } catch (e) {
      console.error(e);
      toast.error("Ошибка", { id: t });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="border border-[var(--border-color)] rounded-[12px] overflow-hidden bg-[var(--bg-primary)] transition hover:shadow-md">

      {/* header */}
      <div className="px-[16px] py-[12px] border-b border-[var(--border-color)] flex items-center justify-between">
        <span className="font-medium text-[var(--text-primary)]">
          Смена №{task.shift}
        </span>

        <span
          className={`text-[12px] px-[8px] py-[2px] rounded-[6px]
            ${
              isStarted
                ? "bg-[var(--status-success-bg)] text-[var(--status-success-text)]"
                : "bg-[var(--status-warning-bg)] text-[var(--status-warning-text)]"
            }`}
        >
          {isStarted ? "Активно" : "Ожидание"}
        </span>
      </div>

      {/* body */}
      <div className="p-[16px] space-y-[6px] text-[14px]">
        <div className="text-[var(--text-primary)]">
          <b>Дата:</b> {date.toLocaleDateString()}
        </div>

        <div className="text-[var(--text-primary)]">
          <b>Время:</b> {date.toLocaleTimeString()}
        </div>
      </div>

      {/* footer */}
      <div className="p-[12px] border-t border-[var(--border-color)]">
        <button
          disabled={!isStarted || loading}
          onClick={handleTake}
          className={`w-full py-[10px] rounded-[10px] font-medium transition
            ${
              isStarted
                ? "bg-[var(--color-primary)] text-white hover:bg-[var(--color-primary-hover)]"
                : "bg-[var(--bg-disabled)] text-[var(--text-secondary)] cursor-not-allowed"
            }
          `}
        >
          {loading ? "Загрузка..." : isStarted ? "Взять задание" : "Начнётся позже"}
        </button>
      </div>
    </div>
  );
};

export default TaskCard;