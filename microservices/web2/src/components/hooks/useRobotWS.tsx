import { useEffect, useReducer, useRef } from "react";
import type { TaskRecord } from "../../types/task";

// ================= TYPES =================
export type OperatorAction = "RETRY" | "SKIP" | "ABORT";

export type MachineState = {
  connection: "connecting" | "connected" | "disconnected";
  deviceAlive: boolean;
  deviceReady: boolean;
  status: string;
  beginState: "idle" | "running" | "error" | "done";
  iteration: number | null;
  availableActions: OperatorAction[] | null;
  error: string | null;
  deviceError: "DISCONNECTED" | "PORT_ERROR" | "BUSY" | "CALIBRATION" | null;
};

type Action =
  | { type: "WS_OPEN" }
  | { type: "WS_CLOSE" }
  | { type: "ACK BOOT" }
  | { type: "STATUS"; status: string; message?: string }
  | {
      type: "UPDATE_STATE";
      status: string;
      iteration?: number | null;
      error?: string | null;
      actions?: OperatorAction[] | null;
      isCriticalError?: boolean;
    }
  | { type: "ACTIONS"; actions: OperatorAction[] }
  | { type: "RESET_ERROR" }
  | { type: "DEVICE_STATUS"; connected: boolean; message?: string };

// ================= MAPPER =================
const mapStateToBegin = (status: string): MachineState["beginState"] => {
  if (["WAIT_READY", "STAND BY"].includes(status)) return "idle";
  if (["BEGIN_ACK"].includes(status)) return "running";
  if (["DONE", "RETURN_DONE"].includes(status)) return "done";
  if (["WAIT_ACTION"].includes(status)) return "error";
  return "running";
};

// ================= REDUCER =================
export const Reducer = (state: MachineState, action: Action): MachineState => {
  switch (action.type) {
    case "WS_OPEN":
      return { ...state, connection: "connected" };

    case "WS_CLOSE":
      return {
        ...state,
        connection: "disconnected",
        deviceAlive: false,
        deviceReady: false,
      };

    case "ACK BOOT":
      return { ...state, deviceAlive: true };

    case "STATUS":
      return {
        ...state,
        deviceReady: action.status === "READY",
        deviceError: action.message?.includes("Port is nil") ? "PORT_ERROR" : state.deviceError,
      };

    case "UPDATE_STATE": {
      const incomingError = action.error;
      const isCritical = action.isCriticalError || !!incomingError;

      let nextError = state.error;
      let nextBeginState = state.beginState;

      if (isCritical) {
        // 🔥 Жестко фиксируем новую ошибку
        nextError = incomingError ?? "Неизвестная ошибка";
        nextBeginState = "error";
      } else {
        // 🔥 Логика "Липкой ошибки":
        // Если у нас уже есть ошибка, и пришло сообщение БЕЗ ошибки (например END: OK),
        // мы ИГНОРИРУЕМ его и сохраняем текущую ошибку.
        if (state.error) {
          nextError = state.error;
          nextBeginState = "error";
        } else {
          // Чистое состояние, ошибок не было и нет
          nextError = null;
          nextBeginState = mapStateToBegin(action.status);
        }
      }

      return {
        ...state,
        status: action.status,
        iteration: action.iteration ?? state.iteration,
        availableActions: action.actions ?? state.availableActions,
        beginState: nextBeginState,
        error: nextError,
      };
    }

    case "ACTIONS":
      return { ...state, availableActions: action.actions };

    case "RESET_ERROR":
      return { 
        ...state, 
        error: null, 
        deviceError: null, 
        beginState: "idle",
        availableActions: null 
      };

    case "DEVICE_STATUS": {
      if (action.connected) {
        return {
          ...state,
          deviceError: null,
          error: null,
          deviceAlive: true,
          connection: "connected",
          beginState: "idle",
          availableActions: null,
        };
      }
      return {
        ...state,
        deviceError: "DISCONNECTED",
        deviceAlive: false,
        beginState: "error",
        error: action.message ?? "Устройство отключено",
        availableActions: null,
      };
    }

    default:
      return state;
  }
};

// ================= HOOK =================
const WS_URL = "/ws";

export function useWSConnection(token: string | null) {
  const wsRef = useRef<WebSocket | null>(null);
  const isConnectingRef = useRef(false);
  const queueRef = useRef<any[]>([]);

  const [state, dispatch] = useReducer(Reducer, {
    connection: "connecting",
    deviceAlive: false,
    deviceReady: false,
    status: "STAND BY",
    beginState: "idle",
    iteration: null,
    availableActions: null,
    error: null,
    deviceError: null,
  } as MachineState);

  const sendMessage = (msg: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    } else {
      queueRef.current.push(msg);
    }
  };

  const startBegin = (record: TaskRecord) => {
    if (state.deviceError === "DISCONNECTED") return;
    if (!["idle", "done", "error"].includes(state.beginState)) return;

    dispatch({ type: "RESET_ERROR" });

    sendMessage({
      type: "START",
      params: {
        shift: record.shift,
        number: record.number,
        receipt: record.receipt,
        required_amount: record.required_amount,
        bunker: record.bunker,
        gcode: record.gcode,
        extraMode: false,
        seed: record.seed,
      },
    });
  };

  const sendOperatorAction = (choice: OperatorAction) => {
    if (state.deviceError === "DISCONNECTED") return;
    dispatch({ type: "RESET_ERROR" });
    sendMessage({ type: "BEGIN", choice });
  };

  const handleDeviceMessage = (msg: any): boolean => {
    if (msg.type === "DEVICE") {
      if (msg.status === "ERROR" && msg.message === "DISCONNECTED") {
        dispatch({ type: "DEVICE_STATUS", connected: false, message: "Устройство отключено" });
        return true;
      }
      if (msg.status === "OK" && msg.message === "CONNECTED") {
        dispatch({ type: "DEVICE_STATUS", connected: true });
        return true;
      }
    }
    return false;
  };

  const connect = () => {
    if (wsRef.current || isConnectingRef.current) return;
    isConnectingRef.current = true;

    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      isConnectingRef.current = false;
      dispatch({ type: "WS_OPEN" });
      sendMessage({ type: "AUTH", token: "Bearer " + token });
      queueRef.current.forEach((m) => ws.send(JSON.stringify(m)));
      queueRef.current = [];
    };

    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data);
      console.log("WS:", msg);

      if (handleDeviceMessage(msg)) return;

      if (msg.type === "AUTH" && msg.status === "OK") {
        sendMessage({ type: "BOOT" });
        sendMessage({ type: "STATUS" });
        return;
      }

      if (msg.type === "BOOT") {
        if (msg.status?.includes("ACK")) dispatch({ type: "ACK BOOT" });
        return;
      }

      if (msg.type === "STATUS") {
        dispatch({ type: "STATUS", status: msg.message ?? msg.status, message: msg.message });
        return;
      }

      // 🔥 ОБЪЕДИНЕННАЯ ОБРАБОТКА: STATE, ACTION, RETURN, BEGIN, END
      const isStatusMessage = ["STATE", "ACTION", "RETURN", "BEGIN", "END"].includes(msg.type);
      
      if (isStatusMessage) {
        const isErrorStatus = msg.status === "ERROR";
        
        // 🔥 Если пришла ошибка, мы её обрабатываем.
        // Если ошибки нет, но у нас уже висит state.error, редюсер сам решит (игнорировать или обновить).
        
        dispatch({
          type: "UPDATE_STATE",
          // Для END: OK статус будет "OK" или сообщение. 
          // Но так как isErrorStatus=false, сработает логика "липкой ошибки" в редюсере.
          status: isErrorStatus ? "ERROR" : (msg.message || msg.status || "UNKNOWN"),
          iteration: msg.Iteration ?? msg.iteration ?? null,
          error: isErrorStatus ? (msg.error ?? msg.message ?? "Ошибка выполнения") : null,
          actions: (msg.actions ?? msg.availableActions ?? msg.choices) as OperatorAction[] | null,
          isCriticalError: isErrorStatus,
        });
        return;
      }

      if (msg.actions) {
        dispatch({ type: "ACTIONS", actions: msg.actions as OperatorAction[] });
      }
    };

    ws.onclose = () => {
      isConnectingRef.current = false;
      dispatch({ type: "WS_CLOSE" });
      wsRef.current = null;
      setTimeout(connect, 2000);
    };

    ws.onerror = () => ws.close();
  };

  useEffect(() => {
    connect();
    return () => { wsRef.current?.close(); };
  }, []);

  const isConnected = state.connection === "connected" && state.deviceAlive && state.deviceError !== "DISCONNECTED";
  const isFullyDisabled = state.deviceError === "DISCONNECTED" || state.connection === "disconnected";

  return {
    sendMessage, startBegin, sendOperatorAction,
    rawStatus: state.status, beginState: state.beginState,
    availableActions: state.availableActions, error: state.error,
    deviceError: state.deviceError, iteration: state.iteration,
    connectionStatus: state.connection, isConnected, isFullyDisabled,
  };
}