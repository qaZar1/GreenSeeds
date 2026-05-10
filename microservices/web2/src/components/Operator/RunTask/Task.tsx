import React from "react";
import type { TaskRecord } from "../../../types/task";
import { useWSConnection } from "../../hooks/useRobotWS";
import { useAuth } from "../../../context/AuthContext";
import ActionButton from "../../utils/AсtionButton";
import { usePageHeader } from "../../../context/HeaderContext";

type Props = {
  record: TaskRecord;
};

// ================= TYPES =================
type OperatorAction = "RETRY" | "SKIP" | "ABORT";

type ActionConfig = {
  label: string;
  bg: string;
  hint?: string;
};

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
  RETRY: {
    label: "🔄 Повторить",
    bg: "!bg-yellow-500 hover:!bg-yellow-600",
    hint: "Повторить последнюю операцию",
  },

  SKIP: {
    label: "⏭️ Пропустить",
    bg: "!bg-[var(--color-primary)] hover:!bg-[var(--color-primary-hover)]",
    hint: "Пропустить текущий шаг",
  },

  ABORT: {
    label: "❌ Отменить",
    bg: "!bg-red-600 hover:!bg-red-700",
    hint: "Прервать задание",
  },
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
const getHumanStatus = (
  status: string,
  beginState: string
) => {
  if (status === "DONE" || status === "RETURN_DONE") {
    return "Завершено";
  }

  if (beginState === "running") {
    return "Выполнение...";
  }

  if (beginState === "error") {
    if (status === "WAIT_ACTION") {
      return "Требуется действие";
    }

    return "Ошибка";
  }

  return STATE_LABELS[status] || status;
};

const getStatusColor = (beginState: string) => {
  if (beginState === "running") {
    return "text-[var(--status-info-text)]";
  }

  if (beginState === "done") {
    return "text-[var(--status-success-text)]";
  }

  if (beginState === "error") {
    return "text-[var(--status-danger-text)]";
  }

  return "text-[var(--text-secondary)]";
};

// ================= STEPPER =================
const Stepper = ({
  current,
}: {
  current: string;
}) => {
  const index = STEPS.indexOf(current);

  const currentIndex =
    index === -1 ? -1 : index;

  return (
    <div className="flex items-start justify-between gap-[12px]">

      {STEPS.map((step, i) => {

        const isActive = current === step;
        const isPassed = currentIndex > i;

        return (
          <div
            key={step}
            className="flex-1 flex flex-col items-center"
          >
            <div
              className={`
                w-6 h-6
                rounded-full
                border-2
                transition-all

                ${
                  isPassed
                    ? "bg-[var(--status-success-text)] border-[var(--status-success-text)]"
                    : isActive
                      ? "bg-[var(--color-primary)] border-[var(--color-primary)]"
                      : "border-[var(--border-color)]"
                }
              `}
            />

            <div className="text-[10px] mt-1 text-[var(--text-secondary)] text-center px-[4px] break-words">
              {STATE_LABELS[step]}
            </div>
          </div>
        );

      })}

    </div>
  );
};

// ================= COMPONENT =================
const TaskCard: React.FC<Props> = ({
  record,
}) => {
  usePageHeader("Выполнение задания", "Управление процессом выполнения");

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

  const validActions =
    availableActions?.filter(
      (
        a
      ): a is OperatorAction =>
        ["RETRY", "SKIP", "ABORT"].includes(a)
    );

  return (
    <div
      className="
        w-full
        max-w-[720px]
        mx-auto

        bg-[var(--bg-card)]
        border border-[var(--border-color)]
        shadow-sm
        rounded-[16px]

        p-[16px] sm:p-5
        space-y-5
      "
    >

      {/* HEADER */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-[12px]">

        <div>
          <div className="text-lg font-semibold text-[var(--text-primary)]">
            Задание №{record.number}
          </div>

          <div className="text-xs text-[var(--text-secondary)]">
            Смена {record.shift}
          </div>
        </div>

        <div className="flex items-center gap-2">

          <div
            className={`
              w-2.5 h-2.5 rounded-full
              ${
                isConnected
                  ? "bg-[var(--status-success-text)]"
                  : "bg-[var(--status-danger-text)]"
              }
            `}
          />

          <span className="text-xs text-[var(--text-secondary)]">
            {isConnected
              ? "Подключено"
              : "Нет связи"}
          </span>

        </div>
      </div>

      {/* SEED */}
      <div className="bg-[var(--bg-page)] rounded-xl p-3 text-center">

        <div className="text-xs text-[var(--text-secondary)]">
          СЕМЕНА
        </div>

        <div className="font-medium uppercase break-words text-[var(--text-primary)]">
          {record.seed_ru}
        </div>

      </div>

      {/* STEPPER */}
      <div className="overflow-x-auto pb-[4px]">
        <div className="min-w-[520px]">
          <Stepper current={rawStatus} />
        </div>
      </div>

      {/* STATUS */}
      <div className="text-center">

        <div className="text-xs text-[var(--text-secondary)]">
          Текущий статус
        </div>

        <div
          className={`
            text-lg
            font-semibold
            ${getStatusColor(beginState)}
          `}
        >
          {getHumanStatus(
            rawStatus,
            beginState
          )}
        </div>

        {beginState === "running" && (
          <div className="text-[var(--status-info-text)] text-sm mt-1 animate-pulse">
            Идёт выполнение...
          </div>
        )}

        {/* DISCONNECTED */}
        {deviceError === "DISCONNECTED" && (
          <div className="mt-3 p-4 rounded-xl bg-[var(--status-danger-bg)] border-2 border-[var(--status-danger-text)]">

            <div className="flex items-start gap-3">

              <span className="text-2xl">
                🔌
              </span>

              <div className="flex-1">

                <div className="text-[var(--status-danger-text)] font-bold text-lg">
                  Устройство отключено
                </div>

                <div className="text-sm text-[var(--status-danger-text)] mt-1 break-words">
                  {error ||
                    "Потеряна связь с устройством. Проверьте кабель питания и USB-подключение."}
                </div>

                <div className="text-xs text-[var(--status-danger-text)] mt-2 bg-[var(--bg-card)] p-2 rounded">
                  💡 После восстановления подключения интерфейс разблокируется автоматически
                </div>

              </div>

            </div>

          </div>
        )}

        {/* ERROR */}
        {beginState === "error" &&
          error &&
          !deviceError && (
            <div className="mt-3 p-4 rounded-xl bg-[var(--status-danger-bg)] border border-[var(--status-danger-text)]">

              <div className="flex flex-col items-center gap-2 text-center">

                <div>

                  <div className="text-[var(--status-danger-text)] font-semibold">
                    Ошибка выполнения
                  </div>

                  <div className="text-sm text-[var(--status-danger-text)] mt-1 break-words">
                    {error}
                  </div>

                </div>

              </div>

              {validActions &&
                validActions.length > 0 && (
                  <>
                    <div className="text-xs text-[var(--text-secondary)] mt-3 mb-2 text-center">
                      Выберите действие:
                    </div>

                    <div className="flex flex-wrap justify-center gap-2 mt-2">

                      {validActions.map(
                        (
                          action: OperatorAction
                        ) => {

                          const config =
                            ACTION_CONFIG[
                              action
                            ];

                          return (
                            <ActionButton
                              key={action}
                              onClick={() =>
                                sendOperatorAction(
                                  action
                                )
                              }
                              disabled={
                                isFullyDisabled
                              }
                              className={`
                                ${config.bg}
                                w-full sm:w-auto
                              `}
                            >
                              {config.label}
                            </ActionButton>
                          );
                        }
                      )}

                    </div>
                  </>
                )}

              {(!validActions ||
                validActions.length ===
                  0) && (
                <div className="text-xs text-[var(--text-secondary)] mt-2 text-center">
                  Обратитесь к оператору
                  или перезапустите задание
                </div>
              )}

            </div>
          )}

        {beginState === "done" && (
          <div className="text-[var(--status-success-text)] text-sm mt-1">
            Операция завершена
          </div>
        )}

        {iteration !== null && (
          <div className="text-xs text-[var(--text-secondary)] mt-1">
            Итерация: {iteration}
          </div>
        )}

      </div>

      {/* ACTIONS */}
      <div className="flex flex-col sm:flex-row gap-2">

        <ActionButton
          onClick={() => startBegin(record)}
          disabled={
            isFullyDisabled ||
            (beginState !== "idle" &&
              beginState !== "done")
          }
          className="
            flex-1
            !bg-[var(--color-primary)]
            hover:!bg-[var(--color-primary-hover)]
          "
        >
          {deviceError === "DISCONNECTED"
            ? "⏳ Ожидание..."
            : "Начать"}
        </ActionButton>

        <ActionButton
          onClick={() =>
            sendMessage({ type: "STOP" })
          }
          disabled={
            isFullyDisabled ||
            beginState !== "running"
          }
          className="
            flex-1
            !bg-red-500
            hover:!bg-red-600
          "
        >
          Стоп
        </ActionButton>

        <ActionButton
          onClick={() =>
            sendMessage({
              type: "SET STATUS READY",
            })
          }
          disabled={isFullyDisabled}
          className="
            flex-1
            !bg-[var(--bg-secondary)]
            hover:!bg-[var(--bg-hover)]
            !text-[var(--text-primary)]
            text-[13px]
            whitespace-nowrap
          "
        >
          DEV: READY
        </ActionButton>

      </div>

      {/* FOOTER */}
      {isFullyDisabled && (
        <div className="text-center text-xs text-[var(--text-secondary)] pt-2 border-t border-[var(--border-color)] break-words px-[4px]">

          {deviceError ===
          "DISCONNECTED"
            ? "🔌 Все функции заблокированы до восстановления связи"
            : "⚠️ Интерфейс временно недоступен"}

        </div>
      )}

    </div>
  );
};

export default TaskCard;