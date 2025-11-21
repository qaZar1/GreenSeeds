import { useEffect, useRef, useState } from "react";
import { useNotify } from "react-admin";

const WS_URL = "/ws";

const robotStateTranslate = {
  "STAND BY": "Ожидает нажатия кнопки",
  READY: "Готов",
  BUSY: "В процессе выполнения",
  ERR: "Ошибка",
  END: "Задание выполнено. Ожидание контроля качества.",
  RETURN: "Возврат каретки",
  UNKNOWN: "Неизвестно",
  MANUAL_MODE: "Ручной режим",
  WAIT: "Выполняется возврат каретки"
};

export function useRobotWS(params = {}) {
  const { record, onSuccessStep, onErrorStep } = params;

  const wsRef = useRef(null);
  const reconnectTimeout = useRef(null);
  const isMountedRef = useRef(false);

  const [state, setState] = useState("STAND BY");
  const [isBoot, setIsBoot] = useState(false);
  const [dots, setDots] = useState("");
  const [control, setControl] = useState(null);
  const [turn, setTurn] = useState(0);
  const [decisionModal, setDecisionModal] = useState({
    open: false,
    reason: "",
    photo: null,
  });

  const notify = useNotify();

  const handleWSMessage = (msg) => {
    const { type, status, payload, error } = msg;

    if (type === "ERR"){
      notify(error, {type:"error"});
      if (onErrorStep && msg.params?.turn) {
        onErrorStep(msg.params.turn, error);
      }
    
      setIsBoot(false);
      setState("ERR");
      return;
    }

    switch (type) {
      case "BOOT":
        if (status === "ACK BOOT") {
          setIsBoot(true);
          break;
        }

        setIsBoot(false);
        break;

      case "STATUS":
        if (status === "STAND BY" ||
          status === "READY" ||
          status === "BUSY" ||
          status === "WAIT" ||
          status === "RETURN" ||
          status === "ERR" ||
          status === "MANUAL_MODE") {
          setState(status);
          break;
        }
        setState("UNKNOWN");
        break;
      
      case "BEGIN":
        if (msg.status.includes("RETURN")) {
          setState("RETURN");
        } else if (msg.status.includes("END")) {
          setState("END");
        } else if (msg.status.includes("ACK")) {
          setState("BUSY");
          setControl(null);
        } else if (msg.status.includes("STAND BY")) {
          setState("STAND BY");
          if (msg.error !== undefined) {
            notify(msg.error, {type:"error"});
            if (onErrorStep && msg.params?.turn) {
              onErrorStep(msg.params.turn, msg.error);
            }

            if (payload && "control" in payload) {
              setControl(payload.control);
            }
            break;
          } else if (msg.params && typeof msg.params.turn !== "undefined") {
            setTurn(msg.params.turn);
          }
          notify("Задание успешно выполнено!", {type: "success"});
          if (onSuccessStep) onSuccessStep();
        }

        if (payload && "control" in payload) {
          setControl(payload.control);
        }
        break;

      case "NEED_DECISION":
        if (msg.payload && typeof msg.payload.reason !== "undefined") {
          setDecisionModal({
            open: true,
            reason: msg.payload.reason,
            photo: msg.payload.photo,
          });
        }
        break;
      
      case "SETSTATUS READY":
        if (status === "ACK SETSTATUS READY") {
          setState("READY");
          break;
        }
        setState("UNKNOWN");
        break;

      default:
        console.warn("Unknown WS message type:", type);
        break;
    }
  };

  useEffect(() => {
    if (state !== "BUSY") {
      setDots("");
      return;
    }
    const interval = setInterval(() => {
      setDots(prev => (prev.length >= 3 ? "" : prev + "."));
    }, 500);
    return () => clearInterval(interval);
  }, [state]);

  const connectWS = () => {
    if (!isMountedRef.current) return;

    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");
      sendMessage({ type: "BOOT" });
      setTimeout(() => {
        sendMessage({ type: "STATUS" });
      }, 2000);
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        console.log("📩 WS message:", msg);

        handleWSMessage(msg);
      } catch (err) {
        console.error("❌ Failed to parse WS message:", event.data, err);
      }
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
      if (isMountedRef.current) {
        console.log("Reconnecting in 2s...");
        reconnectTimeout.current = setTimeout(connectWS, 2000);
      }
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
      ws.close();
    };
  };

  useEffect(() => {
    isMountedRef.current = true;
    connectWS();

    return () => {
      isMountedRef.current = false;
      clearTimeout(reconnectTimeout.current);
      wsRef.current?.close();
    };
  }, []);

  const displayState =
    state === "BUSY"
      ? robotStateTranslate[state] + dots
      : robotStateTranslate[state] || state;

  const sendMessage = (obj) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(obj));
    } else {
      console.warn("WS not ready, message not sent:", obj);
    }
  };

  return {
    rawState: state,
    displayState,
    sendMessage,
    isBoot,
    control,
    turn,
    decisionModal,
    setDecisionModal
  };
}
