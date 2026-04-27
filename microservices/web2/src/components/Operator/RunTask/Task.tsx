import React from "react";
import type { TaskRecord } from "../../../types/task";
import { useWSConnection } from "../../hooks/useRobotWS";
import { useAuth } from "../../../context/AuthContext";

type Props = {
  record: TaskRecord;
};

// ================= TYPES =================
type OperatorAction = "RETRY" | "SKIP" | "ABORT";
type ActionConfig = { label: string; bg: string; hint?: string };

// ================= LABELS =================
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

// ================= ACTION BUTTON CONFIG =================
const ACTION_CONFIG: Record<OperatorAction, ActionConfig> = {
  RETRY: { label: "🔄 Повторить", bg: "bg-yellow-500", hint: "Повторить последнюю операцию" },
  SKIP:  { label: "⏭️ Пропустить", bg: "bg-blue-500", hint: "Пропустить текущий шаг" },
  ABORT: { label: "❌ Отменить", bg: "bg-red-600", hint: "Прервать задание" },
};

// ================= STEPS =================
const STEPS = [
  "WAIT_READY",
  "BEGIN_ACK",
  "PHOTO_DONE",
  "AI_OK",
  "DONE",
];

// ================= HELPERS =================
const getHumanStatus = (status: string, beginState: string) => {
  if (status === "DONE" || status === "RETURN_DONE") return "Завершено";
  if (beginState === "running") return "Выполнение...";
  if (beginState === "error") {
    if (status === "WAIT_ACTION") return "Требуется действие";
    return "Ошибка";
  }
  return STATE_LABELS[status] || status;
};

const getStatusColor = (beginState: string) => {
  if (beginState === "running") return "text-blue-600";
  if (beginState === "done") return "text-green-600";
  if (beginState === "error") return "text-red-600";
  return "text-gray-600";
};

// ================= STEPPER =================
const Stepper = ({ current }: { current: string }) => {
  const index = STEPS.indexOf(current);
  const currentIndex = index === -1 ? -1 : index;

  return (
    <div className="flex items-center justify-between">
      {STEPS.map((step, i) => {
        const isActive = current === step;
        const isPassed = currentIndex > i;
        return (
          <div key={step} className="flex-1 flex flex-col items-center">
            <div
              className={`w-6 h-6 rounded-full border-2 transition-all
                ${isPassed ? "bg-green-500 border-green-500" : isActive ? "bg-blue-500 border-blue-500" : "border-gray-300"}`}
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

// ================= COMPONENT =================
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
    iteration,
    error,
    deviceError,
    isFullyDisabled,
  } = useWSConnection(token);

  const validActions = availableActions?.filter(
    (a): a is OperatorAction => ["RETRY", "SKIP", "ABORT"].includes(a)
  );

  return (
    <div className="max-w-[720px] mx-auto bg-white shadow-lg rounded-2xl p-5 space-y-5 border">

      {/* HEADER */}
      <div className="flex justify-between items-center">
        <div>
          <div className="text-lg font-semibold">Задание №{record.number}</div>
          <div className="text-xs text-gray-500">Смена {record.shift}</div>
        </div>
        <div className="flex items-center gap-2">
          <div className={`w-2.5 h-2.5 rounded-full ${isConnected ? "bg-green-500" : "bg-red-500"}`} />
          <span className="text-xs text-gray-500">
            {isConnected ? "Подключено" : "Нет связи"}
          </span>
        </div>
      </div>

      {/* SEED */}
      <div className="bg-gray-50 rounded-xl p-3 text-center">
        <div className="text-xs text-gray-500">СЕМЕНА</div>
        <div className="font-medium uppercase">{record.seed_ru}</div>
      </div>

      {/* STEPPER */}
      <Stepper current={rawStatus} />

      {/* STATUS */}
      <div className="text-center">
        <div className="text-xs text-gray-500">Текущий статус</div>
        <div className={`text-lg font-semibold ${getStatusColor(beginState)}`}>
          {getHumanStatus(rawStatus, beginState)}
        </div>

        {beginState === "running" && (
          <div className="text-blue-500 text-sm mt-1 animate-pulse">Идёт выполнение...</div>
        )}

        {/* 🔌 DISCONNECTED */}
        {deviceError === "DISCONNECTED" && (
          <div className="mt-3 p-4 rounded-xl bg-red-100 border-2 border-red-400">
            <div className="flex items-start gap-3">
              <span className="text-2xl">🔌</span>
              <div className="flex-1">
                <div className="text-red-800 font-bold text-lg">Устройство отключено</div>
                <div className="text-sm text-red-700 mt-1">
                  {error || "Потеряна связь с устройством. Проверьте кабель питания и USB-подключение."}
                </div>
                <div className="text-xs text-red-600 mt-2 bg-red-50 p-2 rounded">
                  💡 После восстановления подключения интерфейс разблокируется автоматически
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Обычные ошибки (не DEVICE) */}
        {beginState === "error" && error && !deviceError && (
          <div className="mt-3 p-4 rounded-xl bg-red-50 border border-red-200">
            <div className="flex flex-col items-center gap-2 text-center">
              <div>
                <div className="text-red-700 font-semibold">Ошибка выполнения</div>
                <div className="text-sm text-red-600 mt-1">{error}</div>
              </div>
            </div>

            {validActions && validActions.length > 0 && (
              <>
                <div className="text-xs text-gray-500 mt-3 mb-2 text-center">
                  Выберите действие:
                </div>
                <div className="flex flex-wrap justify-center gap-2 mt-2">
                  {validActions.map((action: OperatorAction) => {
                    const config = ACTION_CONFIG[action];
                    return (
                      <button
                        key={action}
                        onClick={() => sendOperatorAction(action)}
                        className={`px-4 py-2 rounded-lg text-white text-sm font-medium transition-opacity ${config.bg} hover:opacity-90 disabled:opacity-50`}
                        disabled={isFullyDisabled}
                        title={config.hint}
                      >
                        {config.label}
                      </button>
                    );
                  })}
                </div>
              </>
            )}

            {(!validActions || validActions.length === 0) && (
              <div className="text-xs text-gray-500 mt-2 text-center">
                Обратитесь к оператору или перезапустите задание
              </div>
            )}
          </div>
        )}
        {beginState === "done" && (
          <div className="text-green-500 text-sm mt-1">Операция завершена</div>
        )}

        {iteration !== null && (
          <div className="text-xs text-gray-400 mt-1">Итерация: {iteration}</div>
        )}
      </div>

      {/* ACTIONS */}
      <div className="flex gap-2">
        <button
          disabled={isFullyDisabled || (beginState !== "idle" && beginState !== "done")}
          onClick={() => startBegin(record)}
          className="flex-1 py-2 rounded-lg bg-blue-600 text-white disabled:opacity-40 disabled:cursor-not-allowed"
        >
          {deviceError === "DISCONNECTED" ? "⏳ Ожидание..." : "Начать"}
        </button>
        <button
          disabled={isFullyDisabled || beginState !== "running"}
          onClick={() => sendMessage({ type: "STOP" })}
          className="flex-1 py-2 rounded-lg bg-red-500 text-white disabled:opacity-40 disabled:cursor-not-allowed"
        >
          Стоп
        </button>
        <button
          onClick={() => sendMessage({ type: "SET STATUS READY" })}
          className="flex-1 py-2 rounded-lg bg-gray-500 text-white disabled:opacity-40"
          disabled={isFullyDisabled}
        >
          DEV: READY
        </button>
      </div>

      {/* FOOTER */}
      {isFullyDisabled && (
        <div className="text-center text-xs text-gray-400 pt-2 border-t">
          {deviceError === "DISCONNECTED" 
            ? "🔌 Все функции заблокированы до восстановления связи" 
            : "⚠️ Интерфейс временно недоступен"}
        </div>
      )}
    </div>
  );
};

export default TaskCard;