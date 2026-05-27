import React from "react";
import type { TaskRecord } from "../../../types/task";
import { useWSConnection } from "../../hooks/useRobotWS";
import { useAuth } from "../../../context/AuthContext";
import ActionButton from "../../utils/AсtionButton";
import { usePageHeader } from "../../../context/HeaderContext";
import { Stepper } from "../../utils/Stepper";
import { useEffect, useState } from "react";

type Props = {
  record: TaskRecord;
};

const STEP_LABELS = [
  "Ожидание",
  "Посев",
  "Фото",
  "Контроль",
  "Завершение",
];

const STEP_INDEX: Record<string, number> = {
  WAIT_READY: 0,
  "STAND BY": 0,

  BEGIN: 1,
  BEGIN_START: 1,
  BEGIN_DONE: 1,

  PHOTO: 2,
  PHOTO_START: 2,
  PHOTO_DONE: 2,

  CONTROL: 3,
  CONTROL_START: 3,
  CONTROL_DONE: 3,

  PROCESS: 4,
  PROCESS_START: 4,
  PROCESS_DONE: 4,

  RETURN_DONE: 4,
  DONE: 4,
};

const getHumanStatus = (
  status: string,
  beginState: string,
  dots: string,
) => {

  if (beginState === "error") {
    return "Ошибка";
  }

  if (status === "DONE" || status === "END") {
    return "Завершено";
  }

  if (beginState === "running") {
    return `Выполнение${dots}`;
  }

  if (beginState === "manual") {
    return "Ручной режим";
  }

  return "Ожидание";
};

const getStatusColor = (
  beginState: string,
) => {

  if (beginState === "running") {
    return "text-[var(--status-info-text)]";
  }

  if (beginState === "done") {
    return "text-[var(--status-success-text)]";
  }

  if (beginState === "error") {
    return "text-[var(--status-danger-text)]";
  }

  if (beginState === "manual") {
    return "text-[var(--status-warning-text)]";
  }

  return "text-[var(--text-secondary)]";
};

// ================= COMPONENT =================
const TaskCard: React.FC<Props> = ({
  record,
}) => {

  usePageHeader(
    "Выполнение задания",
    "Управление процессом выполнения",
  );

  if (!record) {
    return null;
  }

  const { token } = useAuth();

  const {
    sendMessage,
    stopProcess,
    rawStatus,
    isConnected,
    beginState,
    startBegin,
    iteration,
    error,
    deviceError,
    isFullyDisabled,
    message,
  } = useWSConnection(token);

  const isManualMode = beginState === "manual";

  const [dots, setDots] = useState("");

  useEffect(() => {

    if (beginState !== "running") {
      setDots("");
      return;
    }

    const interval = setInterval(() => {

      setDots((prev) => {

        if (prev.length >= 3) {
          return "";
        }

        return prev + ".";
      });

    }, 500);

    return () => clearInterval(interval);

  }, [beginState]);

  const hiddenMessages = [
    "WAIT",
    "READY",
    "OK",
  ];

  const shouldShowMessage =
    message &&
    !hiddenMessages.includes(message);

  const canStart =
    isConnected &&
    !isFullyDisabled &&
    ["idle", "error", "done"].includes(beginState);

  const canStop =
    !isFullyDisabled &&
    beginState === "running";

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
      <Stepper
        steps={STEP_LABELS}
        current={STEP_INDEX[rawStatus] ?? 0}
      />

      {/* MANUAL MODE */}
      {isManualMode ? (
        <div className="rounded-xl border border-yellow-500 bg-yellow-500/10 p-4 text-center">

          <div className="text-lg font-semibold text-yellow-500">
            Включен ручной режим
          </div>

          <div className="text-sm text-[var(--text-secondary)] mt-1">
            Устройство не принимает внешние команды
          </div>

        </div>
      ) : (
        <>
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
                beginState,
                dots,
              )}
            </div>

            {shouldShowMessage &&
              beginState !== "error" && (
                <div className="text-[var(--text-secondary)] text-sm mt-1">
                  {message}
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
          <div className="flex flex-col sm:flex-row gap-2">
            <ActionButton
              onClick={() => startBegin(record)}
              disabled={!canStart}
              className="
                flex-1
                !bg-[var(--color-primary)]
                hover:!bg-[var(--color-primary-hover)]
              "
            >
              {deviceError === "DISCONNECTED"
                ? "Ожидание..."
                : "Начать"}
            </ActionButton>

            <ActionButton
              onClick={() =>
                stopProcess()
              }
              disabled={!canStop}
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
          </>
        )}

      {/* FOOTER */}
      {isFullyDisabled && (
        <div className="text-center text-xs text-[var(--text-secondary)] pt-2 border-t border-[var(--border-color)] break-words px-[4px]">
          {deviceError ===
          "DISCONNECTED"
            ? "Все функции заблокированы до восстановления связи"
            : "Интерфейс временно недоступен"}
        </div>
      )}

    </div>
  );
};

export default TaskCard;