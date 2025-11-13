import { useEffect, useRef, useState } from "react";
import { encodeMsg, encodeGcode } from "../utils/robotProtocol";
import { useNotify } from "react-admin";

const WS_URL = "ws://localhost:8001/ws";

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
  const { record } = params;

  const wsRef = useRef(null);
  const reconnectTimeout = useRef(null);
  const isMountedRef = useRef(false);

  const [state, setState] = useState("STAND BY");
  const [isBoot, setIsBoot] = useState(false);
  const [dots, setDots] = useState("");
  const [control, setControl] = useState(false);
  const [amount, setAmount] = useState(0);

  const notify = useNotify();

  const handleWSMessage = (msg) => {
    const { type, status, payload, error } = msg;

    if (error){
      notify(error, {type:"error"});
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
        if (msg.status.includes("ACK")) {
          setState("BUSY");
        } else if (msg.status.includes("END")) {
          setState("END");
        } else if (msg.status.includes("RETURN")) {
          setState("RETURN");
        } else if (msg.status.includes("STAND BY")) {
          if (msg.params && typeof msg.params.amount !== "undefined") {
            setAmount(msg.params.amount);
          }
          setState("STAND BY");
          notify("Задание успешно выполнено!", {type: "success"});
        }

        if (msg.payload && typeof msg.payload.control !== "undefined") {
          setControl(msg.payload.control);
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

  // useEffect(() => {
  //   if (state === "END" && record) {
  //     fetch("/api/reports/add", {
  //       method: "POST",
  //       body: JSON.stringify({
  //         shift: record.shift,
  //         number: record.number,
  //         receipt: record.receipt,
  //         turn: record.turn,
  //         success: true,
  //       }),
  //     })
  //       .then(res => res.json())
  //       .then(data => console.log("END API called:", data))
  //       .catch(err => console.error(err));
  //   }
  // }, [state, record]);

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
    sendGcode: (gcode, displayText, bunker, shift, number, amount) =>
      wsRef.current?.send(
        encodeGcode(gcode, displayText, bunker, shift, number, amount)
      ),
    isBoot,
    control,
    amount,
  };
}
