import { useEffect, useRef, useState } from "react";
import type { TaskRecord } from "../../types/task";
import type { Report } from "../../types/reports";

type ConnectionStatus = "connecting" | "connected" | "disconnected";
type BeginState = "idle" | "running" | "error" | "done";

type UseWSConnectionParams = {
  token: string | null;
  onMessage: (msg: any) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onFatal?: () => void;
};

const WS_URL = "/ws";
const MAX_RECONNECTS = 5;

export function useWSConnection({
  token,
  onMessage,
  onOpen,
  onClose,
  onFatal,
}: UseWSConnectionParams) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  const reconnectAttempts = useRef(0);
  const isMountedRef = useRef(false);

  // 🔥 очередь сообщений
  const queueRef = useRef<any[]>([]);

  const [rawStatus, setRawStatus] = useState<string>("STAND BY");
  const [isConnected, setIsConnected] = useState<boolean>(false);
  const [availableActions, setAvailableActions] = useState<string[] | null>(null);
  const [connectionStatus, setConnectionStatus] =
    useState<ConnectionStatus>("connecting");

  const [beginState, setBeginState] = useState<BeginState>("idle");

  // ================= SEND =================
  const sendMessage = (obj: unknown) => {
    console.log("WS STATE:", wsRef.current?.readyState);
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(obj));
    } else {
      queueRef.current.push(obj);
    }
  };

  // ================= BEGIN =================
  const startBegin = (record: TaskRecord, reports: Report[] = []) => {
    console.log("START BEGIN, status:", rawStatus, isConnected, beginState);

    // 🔒 защита от неправильных состояний
    if (!isConnected) {
      console.warn("Device not connected");
      return;
    }

    if (rawStatus !== "READY" && rawStatus !== "WAIT") {
      console.warn("Device not ready:", rawStatus);
      return;
    }

    if (beginState !== "idle" && beginState !== "done") {
      console.warn("Begin already in progress or error state");
      return;
    }

    const nextTurn =
      reports.length > 0
        ? Math.max(...reports.map((r) => r.turn)) + 1
        : 1;

    sendMessage({
      type: "BEGIN",
      id: record.id,
      params: {
        shift: record.shift,
        number: record.number,
        receipt: record.receipt,
        turn: nextTurn,
        required_amount: record.required_amount,
        bunker: record.bunker,
        gcode: record.gcode,
        extraMode: false,
        seed: record.seed,
      },
    });

    setBeginState("running");
  };

  // ================= OPERATOR =================
  const sendOperatorAction = (action: "RETRY" | "SKIP" | "ABORT") => {
    console.log("SEND ACTION:", action);

    sendMessage({
      type: "BEGIN", // 🔥 ВАЖНО: оставляем BEGIN
      choice: action,
    });
  };

  // ================= CONNECT =================
  const connectWS = () => {
    if (wsRef.current) return;
    if (!isMountedRef.current) return;

    if (reconnectAttempts.current >= MAX_RECONNECTS) {
      onFatal?.();
      return;
    }

    setConnectionStatus("connecting");

    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      setConnectionStatus("connected");
      reconnectAttempts.current = 0;

      sendMessage({ type: "AUTH", token: "Bearer " + token });

      // 🔥 отправка очереди
      queueRef.current.forEach((msg) =>
        ws.send(JSON.stringify(msg))
      );
      queueRef.current = [];

      onOpen?.();
    };

    ws.onmessage = (event: MessageEvent) => {
      try {
        const msg = JSON.parse(event.data);
        console.log("WS:", msg);

        if (msg.type === "AUTH" && msg.status === "OK") {
          sendMessage({ type: "STATUS" });
          sendMessage({ type: "BOOT" });
        }

        switch (msg.type) {
          case "STATUS":
            setRawStatus(msg.status);
            break;

          case "SET STATUS READY":
            sendMessage({ type: "STATUS" });
            break;

          case "DEVICE":
            if (msg.status?.includes("ACK")) {
              setIsConnected(true);
            }
            break;

          case "BOOT":
            if (msg.status === "OK") {
              setIsConnected(true);
            }
            break;

          case "BEGIN":
            setRawStatus(msg.status);

            if (msg.actions) {
              setAvailableActions(msg.actions);
            } else {
              setAvailableActions(null);
            }

            if (msg.status === "START") {
              setBeginState("running");
            } else if (msg.status === "END") {
              setBeginState("done");
            } else if (
              msg.status === "ERROR" ||
              msg.status === "WAIT_ACTION"
            ) {
              setBeginState("error");
            }

            break;
        }

        onMessage(msg);
      } catch (err) {
        console.error("WS parse error:", err);
      }
    };

    ws.onclose = (e) => {
      console.log("WS CLOSED", e);

      wsRef.current = null;
      if (!isMountedRef.current) return;

      setConnectionStatus("disconnected");
      reconnectAttempts.current += 1;

      const delay = Math.min(2000 * reconnectAttempts.current, 10000);
      reconnectTimeout.current = setTimeout(connectWS, delay);

      onClose?.();
    };

    ws.onerror = () => {
      setConnectionStatus("disconnected");
      ws.close();
    };
  };

  useEffect(() => {
    isMountedRef.current = true;
    connectWS();

    return () => {
      isMountedRef.current = false;

      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }

      wsRef.current?.close();
    };
  }, []);

  return {
    sendMessage,
    connectionStatus,
    rawStatus,
    isConnected,

    beginState,
    startBegin,
    sendOperatorAction,
    availableActions,
  };
}