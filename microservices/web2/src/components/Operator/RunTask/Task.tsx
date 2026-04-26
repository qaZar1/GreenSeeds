import React from "react";
import type { TaskRecord } from "../../../types/task";
import { useWSConnection } from "../../hooks/useRobotWS";
import { useAuth } from "../../../context/AuthContext";

type Props = {
  record: TaskRecord;
};

const STATE_LABELS: Record<string, string> = {
  WAIT_READY: "Ожидание готовности",
  BEGIN_ACK: "Запуск",
  BEGIN_END: "Завершение операции",
  PHOTO_DONE: "Фото сделано",
  AI_OK: "Анализ завершён",
  DONE: "Готово",
  RETURN_DONE: "Возврат завершён",
  "STAND BY": "Ожидание",
};

const STEPS = [
  "WAIT_READY",
  "BEGIN_ACK",
  "PHOTO_DONE",
  "AI_OK",
  "DONE",
];

const getStatusColor = (status: string) => {
  if (status === "DONE") return "text-green-600";
  if (status === "ERROR") return "text-red-600";
  if (status === "WAIT_ACTION") return "text-yellow-600";
  return "text-blue-600";
};

const Stepper = ({ current }: { current: string }) => {
  const currentIndex = STEPS.indexOf(current);

  return (
    <div className="flex items-center justify-between">
      {STEPS.map((step, i) => {
        const isActive = current === step;
        const isPassed = currentIndex > i;

        return (
          <div key={step} className="flex-1 flex flex-col items-center">
            <div
              className={`w-6 h-6 rounded-full border-2
                ${
                  isPassed
                    ? "bg-green-500 border-green-500"
                    : isActive
                    ? "bg-blue-500 border-blue-500"
                    : "border-gray-300"
                }`}
            />
            <div className="text-[10px] mt-1 text-gray-500 text-center">
              {STATE_LABELS[step]}
            </div>
          </div>
        );
      })}
    </div>
  );
};

const TaskCard: React.FC<Props> = ({ record }) => {
  if (!record) return null;

  const { token } = useAuth();

  const {
    sendMessage,
    rawStatus,
    isConnected,
    beginState,
    startBegin,
    sendOperatorAction,
    availableActions,
  } = useWSConnection(token);

  return (
    <div className="max-w-[720px] mx-auto bg-white shadow-lg rounded-2xl p-5 space-y-5 border">

      {/* HEADER */}
      <div className="flex justify-between items-center">
        <div>
          <div className="text-lg font-semibold">
            Задание №{record.number}
          </div>
          <div className="text-xs text-gray-500">
            Смена {record.shift}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <div
            className={`w-2.5 h-2.5 rounded-full ${
              isConnected ? "bg-green-500" : "bg-red-500"
            }`}
          />
          <span className="text-xs text-gray-500">
            {isConnected ? "Подключено" : "Нет связи"}
          </span>
        </div>
      </div>

      {/* SEED */}
      <div className="bg-gray-50 rounded-xl p-3 text-center">
        <div className="text-xs text-gray-500">СЕМЕНА</div>
        <div className="font-medium uppercase">
          {record.seed_ru}
        </div>
      </div>

      {/* STEPPER */}
      <Stepper current={rawStatus} />

      {/* STATUS */}
      <div className="text-center">
        <div className="text-xs text-gray-500">Текущий статус</div>
        <div className={`text-lg font-semibold ${getStatusColor(rawStatus)}`}>
          {STATE_LABELS[rawStatus] || rawStatus}
        </div>

        {beginState === "error" && (
          <div className="text-red-500 text-sm mt-1">
            Требуется действие оператора
          </div>
        )}
      </div>

      {/* ACTIONS */}
      <div className="flex gap-2">
        <button
          disabled={!isConnected || beginState === "running"}
          onClick={() => startBegin(record)}
          className="flex-1 py-2 rounded-lg bg-blue-600 text-white disabled:opacity-40"
        >
          Начать
        </button>

        <button
          disabled={!isConnected}
          onClick={() => sendMessage({ type: "STOP" })}
          className="flex-1 py-2 rounded-lg bg-red-500 text-white disabled:opacity-40"
        >
          Стоп
        </button>
      </div>

      {/* OPERATOR ACTIONS */}
      {availableActions && (
        <div className="flex gap-2">
          {availableActions.map((a) => (
            <button
              key={a}
              onClick={() => sendOperatorAction(a as any)}
              className="flex-1 py-2 rounded-lg bg-yellow-500 text-white"
            >
              {a}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default TaskCard;