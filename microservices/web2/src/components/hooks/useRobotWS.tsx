import { useEffect, useReducer, useRef } from "react";
import type { TaskRecord } from "../../types/task";
import { Reducer } from "./machineState";

const WS_URL = "/ws";
const MAX_RECONNECTS = 5;

export function useWSConnection(token: string | null) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttempts = useRef(0);
  const queueRef = useRef<any[]>([]);

  const [state, dispatch] = useReducer(Reducer, {
    connection: "connecting",
    deviceReady: false,
    status: "STAND BY",
    beginState: "idle",
    iteration: null,
    availableActions: null,
  });

  // ================= SEND =================
  const sendMessage = (msg: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    } else {
      queueRef.current.push(msg);
    }
  };

  // ================= BEGIN =================
  const startBegin = (record: TaskRecord) => {
    if (!state.deviceReady) return;
    if (state.beginState !== "idle" && state.beginState !== "done") return;

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

  // ================= ACTION =================
  const sendOperatorAction = (choice: "RETRY" | "SKIP" | "ABORT") => {
    sendMessage({
      type: "START",
      choice,
    });
  };

  // ================= CONNECT =================
  const connect = () => {
    if (wsRef.current) return;
    if (reconnectAttempts.current >= MAX_RECONNECTS) return;

    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      dispatch({ type: "WS_OPEN" });
      reconnectAttempts.current = 0;

      sendMessage({ type: "AUTH", token: "Bearer " + token });

      queueRef.current.forEach((m) => ws.send(JSON.stringify(m)));
      queueRef.current = [];
    };

    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data);
      console.log("WS:", msg);

      if (msg.type === "AUTH" && msg.status === "OK") {
        dispatch({ type: "AUTH_OK" });
        sendMessage({ type: "STATUS" });
        sendMessage({ type: "BOOT" });
        return;
      }

      if (msg.type === "STATE") {
        dispatch({
          type: "STATE",
          status: msg.status,
          iteration: msg.Iteration,
        });
        return;
      }

      if (msg.type === "DEVICE" && msg.status?.includes("ACK")) {
        dispatch({ type: "DEVICE_ACK" });
        return;
      }

      if (msg.type === "BOOT" && msg.status === "OK") {
        dispatch({ type: "BOOT_OK" });
        return;
      }

      if (msg.actions) {
        dispatch({ type: "ACTIONS", actions: msg.actions });
      }
    };

    ws.onclose = () => {
      dispatch({ type: "WS_CLOSE" });
      wsRef.current = null;

      reconnectAttempts.current += 1;
      setTimeout(connect, 2000);
    };

    ws.onerror = () => {
      ws.close();
    };
  };

  useEffect(() => {
    connect();
    return () => wsRef.current?.close();
  }, []);

  return {
    // 🔥 UI-friendly API
    sendMessage,
    startBegin,
    sendOperatorAction,

    rawStatus: state.status,
    beginState: state.beginState,
    availableActions: state.availableActions,

    connectionStatus: state.connection,
    isConnected: state.connection === "connected" && state.deviceReady,
  };
}