export type MachineState = {
  connection: "connecting" | "connected" | "disconnected";

  deviceAlive: boolean;
  deviceReady: boolean;

  status: string;
  beginState: "idle" | "running" | "error" | "done";

  iteration: number | null;
  availableActions: string[] | null;

  error: string | null;
};

type Action =
  | { type: "WS_OPEN" }
  | { type: "WS_CLOSE" }
  | { type: "ACK BOOT" }
  | { type: "STATUS"; status: string; }
  | {
      type: "STATE";
      status: string;
      iteration?: number | null;
      error?: string | null; // 🔥 добавили
      actions?: ("RETRY" | "SKIP" | "ABORT")[] | null; // 🔥 добавили
    }
  | {
      type: "ACTIONS";
      actions: ("RETRY" | "SKIP" | "ABORT")[];
    };

const mapStateToBegin = (
  status: string
): MachineState["beginState"] => {
  if (["WAIT_READY", "STAND BY"].includes(status)) return "idle";
  if (["BEGIN_ACK"].includes(status)) return "running";
  if (["DONE", "RETURN_DONE"].includes(status)) return "done";
  if (["WAIT_ACTION", "ERROR"].includes(status)) return "error";
  return "running";
};

export const Reducer = (
  state: MachineState,
  action: Action
): MachineState => {
  switch (action.type) {
    case "WS_OPEN":
      return { ...state, connection: "connected" };

    case "WS_CLOSE":
      return {
        ...state,
        connection: "disconnected",
        deviceAlive: false,
        deviceReady: false,
        beginState: "idle", // 🔥 Сбрасываем, чтобы старт стал доступен после реконнекта
      };

    case "ACK BOOT":
      return { ...state, deviceAlive: true };

    case "STATUS":
      return {
        ...state,
        deviceReady: action.status === "READY",
      };

    case "STATE": {
      const isErrorMessage = action.status === "ERROR" || action.error != null;
      
      return {
        ...state,
        status: action.status,
        iteration: action.iteration ?? state.iteration,
        availableActions: action.actions ?? state.availableActions,
        
        // 🔥 Явно ставим error-режим, если пришла ошибка
        beginState: isErrorMessage 
          ? "error" 
          : mapStateToBegin(action.status),
        
        // 🔥 Сохраняем текст ошибки (приоритет: action.error > msg.message)
        error: isErrorMessage ? (action.error ?? "Неизвестная ошибка") : null,
      };
    }

    case "ACTIONS":
      return {
        ...state,
        availableActions: action.actions,
      };

    default:
      return state;
  }
};