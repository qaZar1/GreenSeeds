import React, { useEffect, useState, useRef } from "react";
import type { TaskRecord } from "../../../types/task";
import { useRobotWS } from "../../hooks/useRobotWS";
import { useAuth } from "../../../context/AuthContext";
import ActionButton from "../../utils/AсtionButton";
import { usePageHeader } from "../../../context/HeaderContext";
import { Stepper } from "../../utils/Stepper";
import FormModal from "../../utils/FormModal";

type Props = {
  record: TaskRecord;
};

const STEP_LABELS = ["Ожидание", "Посев", "Фото", "Контроль", "Возврат"];

const STEP_INDEX: Record<string, number> = {
  WAIT_READY: 0,
  BEGIN: 1,
  PHOTO: 2,
  CONTROL: 3,
  RETURN: 4,
};

const STEP_TITLES: Record<string, string> = {
  WAIT_READY: "Ожидание готовности",
  BEGIN: "Посев",
  PHOTO: "Получение фото",
  CONTROL: "Контроль качества",
  RETURN: "Возврат каретки",
};

const ERROR_TITLES: Record<string, string> = {
  DEVICE: "Ошибка устройства",
  USER: "Остановлено пользователем",
  INTERNAL: "Внутренняя ошибка",
  AI: "Ошибка анализа",
  CAMERA: "Ошибка камеры",
};

const TaskCard: React.FC<Props> = ({ record }) => {
  usePageHeader("Выполнение задания", "Управление процессом выполнения");
  const { token } = useAuth();

  const {
    startPlanting,
    stopPlanting,
    setReady,
    step,
    message,
    progress,
    iteration,
    error,
    done,
    stopped,
    connection,
    isConnected,
    hasError,
    isRunning,
    logs,
  } = useRobotWS(token);

  const [dots, setDots] = useState("");
  const [isExtraModeModalOpen, setIsExtraModeModalOpen] = useState(false);
  
  const leftColRef = useRef<HTMLDivElement>(null);
  const [rightColHeight, setRightColHeight] = useState<number | string>("auto");

  const effectiveStep = (hasError && error?.stage) ? error.stage : step;
  const currentStep = STEP_INDEX[effectiveStep || "WAIT_READY"] ?? 0;
  const isDisconnected = connection !== "connected";
  
  const canStop = isConnected && isRunning && !stopped;
  const canStart = isConnected && !isRunning && !stopped;

  useEffect(() => {
    if (!isRunning) {
      setDots("");
      return;
    }
    const interval = setInterval(() => {
      setDots((prev) => (prev.length >= 3 ? "" : prev + "."));
    }, 500);
    return () => clearInterval(interval);
  }, [isRunning]);

  useEffect(() => {
    const updateHeight = () => {
      if (leftColRef.current) {
        setRightColHeight(leftColRef.current.offsetHeight);
      }
    };

    updateHeight();
    
    const observer = new ResizeObserver(updateHeight);
    if (leftColRef.current) {
      observer.observe(leftColRef.current);
    }

    window.addEventListener("resize", updateHeight);

    return () => {
      window.removeEventListener("resize", updateHeight);
      observer.disconnect();
    };
  }, [record, logs.length, effectiveStep, done, stopped, hasError, message]);

  const getStatusText = () => {
    if (done) return "Завершено";
    if (stopped) return "Остановлено";
    if (hasError && error?.stage) return `${STEP_TITLES[error.stage]}: Ошибка`;
    if (effectiveStep) return `${STEP_TITLES[effectiveStep]}${dots}`;
    return "Ожидание";
  };

  const getStatusColor = () => {
    if (done) return "text-[var(--status-success-text)]";
    if (stopped) return "text-yellow-500";
    return "text-[var(--status-info-text)]";
  };

  return (
  <>
    <div className="max-w-[1400px] mx-auto">
      <div className="grid grid-cols-1 xl:grid-cols-[2fr_1fr] gap-4 items-start">

        <div 
          ref={leftColRef}
          className="bg-[var(--bg-card)] border border-[var(--border-color)] shadow-sm rounded-[16px] p-4 sm:p-5 space-y-5"
        >
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <div>
              <div className="text-lg font-semibold text-[var(--text-primary)]">
                Задание №{record.number}
              </div>
              <div className="text-xs text-[var(--text-secondary)]">
                Смена {record.shift}
              </div>
            </div>
            <div className="flex items-center gap-2">
              <div className={`w-2.5 h-2.5 rounded-full ${isConnected ? "bg-[var(--status-success-text)]" : "bg-[var(--status-danger-text)]"}`} />
              <span className="text-xs text-[var(--text-secondary)]">
                {isConnected ? "Подключено" : "Нет связи"}
              </span>
            </div>
          </div>

          {/* Seed Info */}
          <div className="bg-[var(--bg-page)] rounded-xl p-3 text-center">
            <div className="text-xs text-[var(--text-secondary)] uppercase">Семена</div>
            <div className="font-medium text-[var(--text-primary)] break-words">{record.seed_ru}</div>
          </div>

          {/* Stepper */}
          <Stepper 
            steps={STEP_LABELS} 
            current={currentStep} 
            completed={done} 
            errorStepIndex={hasError && error?.stage ? STEP_INDEX[error.stage] : null}
          />

          {/* Main Content */}
          <div className="text-center space-y-4">
            {hasError && (
              <div className="p-4 rounded-xl bg-[var(--status-danger-bg)] border border-[var(--status-danger-text)]">
                <div className="font-bold text-[var(--status-danger-text)]">
                  {ERROR_TITLES[error?.code || ""] || "Ошибка"}
                </div>
                <div className="text-sm mt-1 text-[var(--status-danger-text)]">
                  {error?.message}
                </div>
              </div>
            )}
            
            {!hasError && (
              <div>
                <div className="text-xs text-[var(--text-secondary)] mb-1">Текущий статус</div>
                <div className={`text-lg font-bold ${getStatusColor()}`}>{getStatusText()}</div>
              </div>
            )}

            {progress && (
              <div className="mt-2">
                <div className="flex justify-between text-xs text-[var(--text-secondary)] mb-1">
                  <span>Итерация {progress.current} / {progress.total}</span>
                  <span>{progress.percent}%</span>
                </div>
                <div className="h-2 bg-[var(--bg-page)] rounded-full overflow-hidden">
                  <div className="h-full bg-[var(--color-primary)] transition-all duration-300" style={{ width: `${progress.percent}%` }} />
                </div>
              </div>
            )}

            {!progress && iteration !== null && (
              <div className="text-xs text-[var(--text-secondary)]">Итерация: {iteration}</div>
            )}
          </div>

          {/* Actions */}
          <div className="flex flex-col sm:flex-row gap-2 pt-2">
            <ActionButton onClick={() => setIsExtraModeModalOpen(true)} disabled={!canStart} className="flex-1 !bg-[var(--color-primary)] hover:!bg-[var(--color-primary-hover)]">
              {isDisconnected ? "Нет связи" : "Начать"}
            </ActionButton>
            <ActionButton onClick={stopPlanting} disabled={!canStop} className="flex-1 !bg-red-500 hover:!bg-red-600">Стоп</ActionButton>
            <ActionButton onClick={setReady} disabled={!isConnected} className="flex-1 !bg-[var(--bg-secondary)] hover:!bg-[var(--bg-hover)] !text-[var(--text-primary)] text-[13px] whitespace-nowrap">DEV: READY</ActionButton>
          </div>
        </div>

        <div 
          className="bg-[var(--bg-card)] border border-[var(--border-color)] shadow-sm rounded-[16px] p-4 sm:p-5 flex flex-col overflow-hidden"
          style={{ height: rightColHeight }}
        >
          <div className="flex items-center justify-between mb-4 shrink-0">
            <h3 className="font-semibold text-[var(--text-primary)]">Журнал событий</h3>
            <span className="text-xs text-[var(--text-secondary)]">{logs.length} записей</span>
          </div>
          
          <div className="flex-1 overflow-y-auto min-h-0 pr-2 space-y-2">
            {logs.length === 0 ? (
              <div className="h-full flex items-center justify-center text-sm text-[var(--text-secondary)] opacity-70">Нет событий</div>
            ) : (
              logs
                .slice()
                .reverse()
                .map((log, index) => {
                  const uniqueKey = `${log.id}-${index}`;
                  
                  const isError = log.event === "ERROR";
                  const isDone = log.event === "DONE";
                  const isStop = log.event === "STOP";

                  return (
                    <div key={uniqueKey} className="p-3 rounded-xl bg-[var(--bg-page)] border border-[var(--border-color)]">
                      <div className="flex items-start gap-3">
                        <div className={`mt-1.5 w-2 h-2 rounded-full shrink-0 ${isError ? "bg-red-500" : isDone ? "bg-green-500" : isStop ? "bg-yellow-500" : "bg-[var(--color-primary)]"}`} />
                        <div className="flex-1 min-w-0">
                          <div className={`text-sm ${isError ? "text-red-500" : "text-[var(--text-primary)]"}`}>{log.message}</div>
                          <div className="text-xs text-[var(--text-secondary)] mt-1">
                            {log.time}
                            {log.step && <span className="ml-2">• {STEP_TITLES[log.step]}</span>}
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })
            )}
          </div>
        </div>

      </div>
    </div>

    <FormModal
      isOpen={isExtraModeModalOpen}
      title="Запуск задания"
      onClose={() => setIsExtraModeModalOpen(false)}
      onSubmit={(data) => {
        setIsExtraModeModalOpen(false);
        startPlanting({
          ...record,
          extraMode: data.extraMode === true || data.extraMode === "true",
        });
      }}
      fields={[
        {
          name: "extraMode",
          label: "Дополнительный режим",
          type: "select",
          required: true,
          options: [
            { label: "Выключен", value: false },
            { label: "Включен", value: true },
          ],
        },
      ]}
      initialValues={{ extraMode: false }}
    />
  </>
);
};

export default TaskCard;
