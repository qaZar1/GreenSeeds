import React from "react";
import { useTaskController } from "./useTaskContoller";
import type { TaskRecord } from "../../../types/task";
import { useWSConnection } from "../../hooks/useRobotWS";
import { useAuth } from "../../../context/AuthContext";
import { useState } from "react";

type Props = {
  record: TaskRecord;
};

const TaskCard: React.FC<Props> = ({ record }) => {
  if (!record) return null;

  const {token} = useAuth();

  // ✅ сначала WS
  const {
    sendMessage,
    connectionStatus,
    rawStatus,
    isConnected,
    beginState,
    startBegin,
    sendOperatorAction,
    availableActions} = useWSConnection({
    token,
    onMessage: () => {},
    onOpen: () => {},
    onClose: () => {},
    onFatal: () => {
      alert("Failed to connect to WS");
    },
  });

  // ✅ потом контроллер
  // const ctrl = useTaskController(record, ws);

  // const canStart =
  //   ws.rawState === "READY" &&
  //   !ctrl.isRunning &&
  //   ctrl.completedAmount < ctrl.requiredAmount;

  return (
    <div className="max-w-[700px] mx-auto bg-[var(--bg-card)] border border-[var(--border-color)] rounded-[16px] p-[20px] space-y-[16px]">

      {/* HEADER */}
      <div className="flex justify-between items-center">
        <div className="text-[20px] font-semibold">
          Задание №{record.number}
        </div>

        <div
          className={`w-[10px] h-[10px] rounded-full ${
            isConnected ? "bg-green-500" : "bg-red-500"
          }`}
        />
      </div>

      {/* SEED */}
      <div className="p-[12px] border rounded-[10px] text-center">
        <div className="text-[12px]">СЕМЕНА</div>
        <div className="text-[14px] uppercase">{record.seed_ru}</div>
      </div>

      {/* INFO */}
      <div className="grid grid-cols-2 gap-[12px]">
        <div>
          <div className="text-[12px]">СМЕНА</div>
          <div>{record.shift}</div>
        </div>

        <div className="text-right">
          <div className="text-[12px]">КОНТРОЛЬ</div>
          {/* <div>
            {ws.control === null
              ? "Ожидание"
              : ws.control
              ? "Пройден"
              : "Не пройден"}
          </div> */}
        </div>
      </div>

      {/* PROGRESS */}
      <div>
        {/* <div className="text-[12px] text-center">
          {ctrl.completedAmount} / {ctrl.requiredAmount}
        </div>

        <div className="h-[8px] bg-gray-200 rounded-full">
          <div
            className="h-full bg-blue-500"
            style={{ width: `${ctrl.progress}%` }}
          />
        </div> */}
      </div>

      {/* STATUS */}
      <div className="p-[12px] border rounded-[10px] text-center">
        <div className="font-semibold">{rawStatus}</div>

        {/* {ctrl.isRunning && (
          <div className="text-[12px]">
            Лоток {ws.turn} / {ctrl.requiredAmount}
          </div>
        )}

        {ctrl.isRetrying && (
          <div className="text-red-500 text-[12px]">
            ⚠️ Повтор
          </div>
        )} */}
      </div>

      {/* ACTIONS */}
      <div className="flex gap-[8px]">
        <button
          disabled={!isConnected}
          onClick={() => {
            // ctrl.setIsRunning(true);
            console.log(record);
            startBegin(record, record.reports ?? []);
          }}
          className="flex-1 py-[10px] rounded bg-blue-500 text-white disabled:opacity-50"
        >
          Начать
        </button>

        <button
          disabled={!isConnected}
          onClick={() => {
            // ctrl.setIsRunning(false);
            sendMessage({ type: "STOP" });
          }}
          className="flex-1 py-[10px] rounded bg-red-500 text-white disabled:opacity-50"
        >
          Стоп
        </button>

        <button
          disabled={!isConnected}
          onClick={() => sendMessage({ type: "SET STATUS READY" })}
          className="flex-1 py-[10px] rounded bg-gray-500 text-white"
        >
          DEV: READY
        </button>
      </div>
      {availableActions && (
        <div style={{ marginTop: 20 }}>
          {availableActions.includes("RETRY") && (
            <button
            className="flex-1 py-[10px] rounded bg-gray-500 text-white"
            onClick={() => sendOperatorAction("RETRY")}
            >
              Retry
            </button>
          )}

          {availableActions.includes("SKIP") && (
            <button
            onClick={() => sendOperatorAction("SKIP")}
            className="flex-1 py-[10px] rounded bg-gray-500 text-white"
            >
              Skip
            </button>
          )}

          {availableActions.includes("ABORT") && (
            <button
            onClick={() => sendOperatorAction("ABORT")}
            className="flex-1 py-[10px] rounded bg-gray-500 text-white"
            >
              Abort
            </button>
          )}
        </div>
      )}
    </div>
  );
};

export default TaskCard;