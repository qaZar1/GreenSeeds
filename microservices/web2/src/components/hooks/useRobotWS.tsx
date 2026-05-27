import { useEffect, useReducer, useRef } from "react";
import type { TaskRecord } from "../../types/task";

// ================= TYPES =================
export type MachineState = {
  connection: "connecting" | "connected" | "disconnected";

  deviceAlive: boolean;

  deviceReady: boolean;

  status: string;

  message: string | null;

  beginState: "idle" | "running" | "error" | "done" | "manual";

  iteration: number | null;

  error: string | null;

  deviceError:
    | "DISCONNECTED"
    | "PORT_ERROR"
    | "BUSY"
    | "CALIBRATION"
    | null;
};

type Action =
  | { type: "WS_OPEN" }
  | { type: "WS_CLOSE" }
  | { type: "ACK_BOOT" }
  | {
      type: "STATUS";
      status: string;
      message?: string;
    }
  | {
      type: "UPDATE_STATE";
      status: string;
      message?: string | null;
      iteration?: number | null;
      error?: string | null;
    }
  | { type: "RESET_ERROR" }
  | {
      type: "DEVICE_STATUS";
      connected: boolean;
      message?: string;
    }
  | {
      type: "FORCE_IDLE";
      message?: string | null;
    };

// ================= HELPERS =================
const mapStateToBegin = (
  status: string,
): MachineState["beginState"] => {

  if (status === "WAIT_READY") {
    return "idle";
  }

  if (status === "DONE") {
    return "done";
  }

  if (status === "ERROR") {
    return "error";
  }
  
  if (status === "MANUAL") {
    return "manual";
  }

  return "running";
};

// ================= REDUCER =================
export const Reducer = (
  state: MachineState,
  action: Action,
): MachineState => {

  switch (action.type) {

    case "WS_OPEN":
      return {
        ...state,
        connection: "connected",
      };

    case "WS_CLOSE":
      return {
        ...state,

        connection: "disconnected",

        deviceAlive: false,

        deviceReady: false,
      };

    case "ACK_BOOT":
      return {
        ...state,
        deviceAlive: true,
      };

    case "STATUS":
      return {
        ...state,

        deviceReady: action.status === "READY",

        message:
          action.message ??
          state.message,

        deviceError:
          action.message?.includes("Port is nil")
            ? "PORT_ERROR"
            : state.deviceError,
      };

    case "UPDATE_STATE": {
      const isError =
        !!action.error;

      return {
        ...state,

        status: action.status,

        message:
          action.message ?? state.message,

        iteration:
          action.iteration ?? state.iteration,

        beginState:
          isError
            ? "error"
            : mapStateToBegin(action.status),

        error:
          action.error ?? null,
      };
    }

    case "FORCE_IDLE":
      return {
        ...state,

        status: "WAIT_READY",

        message: action.message ?? null,

        beginState: "idle",

        iteration: null,

        error: null,
      };

    case "RESET_ERROR":
      return {
        ...state,
        error: null,
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
        };
      }

      return {
        ...state,

        deviceError: "DISCONNECTED",

        deviceAlive: false,

        beginState: "error",

        error:
          action.message ??
          "Устройство отключено",
      };
    }

    default:
      return state;
  }
};

// ================= HOOK =================
const WS_URL = "/ws";

export function useWSConnection(
  token: string | null,
) {

  const wsRef =
    useRef<WebSocket | null>(null);

  const isConnectingRef =
    useRef(false);

  const queueRef = useRef<any[]>([]);

  const [state, dispatch] = useReducer(
    Reducer,
    {
      connection: "connecting",

      deviceAlive: false,

      deviceReady: false,

      status: "WAIT_READY",

      message: "Ожидание готовности",

      beginState: "idle",

      iteration: null,

      error: null,

      deviceError: null,
    } as MachineState,
  );

  const sendMessage = (msg: any) => {

    console.log(
      "%c[WS OUT]",
      "color:#4ade80;font-weight:bold",
      msg,
    );

    if (
      wsRef.current?.readyState ===
      WebSocket.OPEN
    ) {

      wsRef.current.send(
        JSON.stringify(msg),
      );

      return;
    }

    console.log(
      "%c[WS QUEUED]",
      "color:#facc15;font-weight:bold",
      msg,
    );

    queueRef.current.push(msg);
  };

  const startBegin = (
    record: TaskRecord,
  ) => {

    if (
      state.deviceError ===
      "DISCONNECTED"
    ) {
      return;
    }

    if (
      ![
        "idle",
        "done",
        "error",
      ].includes(state.beginState)
    ) {
      return;
    }

    dispatch({
      type: "RESET_ERROR",
    });

    sendMessage({
      type: "START",

      params: {
        shift: record.shift,

        number: record.number,

        receipt: record.receipt,

        required_amount:
          record.required_amount,

        bunker: record.bunker,

        gcode: record.gcode,

        extraMode: false,

        seed: record.seed,
      },
    });
  };

  const stopProcess = () => {

    sendMessage({
      type: "STOP",
    });
  };

  const handleDeviceMessage = (
    msg: any,
  ): boolean => {

    if (msg.type !== "DEVICE") {
      return false;
    }

    if (
      msg.status === "ERROR" &&
      msg.message === "DISCONNECTED"
    ) {

      dispatch({
        type: "DEVICE_STATUS",

        connected: false,

        message:
          "Устройство отключено",
      });

      return true;
    }

    if (
      msg.status === "OK" &&
      msg.message === "CONNECTED"
    ) {

      dispatch({
        type: "DEVICE_STATUS",

        connected: true,
      });

      return true;
    }

    return false;
  };

  const connect = () => {

    if (
      wsRef.current ||
      isConnectingRef.current
    ) {
      return;
    }

    isConnectingRef.current = true;

    const ws = new WebSocket(WS_URL);

    wsRef.current = ws;

    ws.onopen = () => {

      isConnectingRef.current = false;

      dispatch({
        type: "WS_OPEN",
      });

      sendMessage({
        type: "AUTH",

        token:
          "Bearer " + token,
      });

      queueRef.current.forEach((m) => {
        ws.send(JSON.stringify(m));
      });

      queueRef.current = [];
    };

    ws.onmessage = (e) => {

      console.log(
        "%c[WS RAW IN]",
        "color:orange;font-weight:bold",
        e.data,
      );

      let msg;

      try {
        msg = JSON.parse(e.data);
      } catch (err) {

        console.error(
          "[WS PARSE ERROR]",
          err,
        );

        return;
      }

      console.log(
        "%c[WS IN]",
        "color:#60a5fa;font-weight:bold",
        msg,
      );

      if (
        handleDeviceMessage(msg)
      ) {
        return;
      }

      if (
        msg.type === "AUTH" &&
        msg.status === "OK"
      ) {

        sendMessage({
          type: "BOOT",
        });

        sendMessage({
          type: "STATUS",
        });

        return;
      }

      if (msg.type === "BOOT") {
        if (msg.status?.includes("ACK")) {
          dispatch({
            type: "ACK_BOOT",
          });
        }

        return;
      }

      if (msg.type === "STATUS") {
        dispatch({
          type: "STATUS",
          status:
            msg.status ??
            "UNKNOWN",
          message:
            msg.message,
        });

        return;
      }

      if (msg.type === "END") {

        dispatch({
          type: "FORCE_IDLE",

          message:
            msg.message ??
            "Остановлено",
        });

        return;
      }

      // backend state messages
      const isProcessMessage = [
        "STATE",
        "BEGIN",
        "PHOTO",
        "CONTROL",
        "PROCESS",
        "END",
      ].includes(msg.type);

      if (isProcessMessage) {

        dispatch({
          type: "UPDATE_STATE",

          status:
            msg.status ??
            "UNKNOWN",

          message:
            msg.message ??
            null,

          iteration:
            msg.iteration ??
            msg.Iteration ??
            null,

          error:
            msg.status === "ERROR"
              ? (
                  msg.message ??
                  "Ошибка выполнения"
                )
              : null,
        });

        return;
      }
    };

    ws.onclose = () => {

      isConnectingRef.current = false;

      dispatch({
        type: "WS_CLOSE",
      });

      wsRef.current = null;

      setTimeout(
        connect,
        2000,
      );
    };

    ws.onerror = () => {
      ws.close();
    };
  };

  useEffect(() => {

    connect();

    return () => {
      wsRef.current?.close();
    };

  }, []);

  const isConnected =
    state.connection ===
      "connected" &&
    state.deviceAlive &&
    state.deviceError !==
      "DISCONNECTED";

  const isFullyDisabled =
    state.deviceError ===
      "DISCONNECTED" ||
    state.connection ===
      "disconnected";

  return {
    sendMessage,

    startBegin,

    stopProcess,

    rawStatus:
      state.status,

    message:
      state.message,

    beginState:
      state.beginState,

    error:
      state.error,

    deviceError:
      state.deviceError,

    iteration:
      state.iteration,

    connectionStatus:
      state.connection,

    isConnected,

    isFullyDisabled,
  };
}